package pubsub

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pinmonl/pinmonl/model"
	"github.com/sirupsen/logrus"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var defaultUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	// The Pubsuber.
	hub Pubsuber

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan Message

	// The subscribed topics of client.
	subs map[string]bool

	user *model.User
}

func NewClient(hub Pubsuber, conn *websocket.Conn, user *model.User) *Client {
	return &Client{
		hub:  hub,
		conn: conn,
		send: make(chan Message, 256),
		subs: make(map[string]bool),
		user: user,
	}
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.Unregister(c)
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var body struct {
			Topic string `json:"topic"`
		}
		err := c.conn.ReadJSON(&body)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.Debugf("pubsub: client read err(%v)", err)
			}
			break
		}

		c.Subscribe(body.Topic)
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.conn.WriteJSON(map[string]interface{}{
				"topic": msg.Topic(),
				"data":  msg.Data(),
			})
			if err != nil {
				return
			}

			// Add queued chat messages to the current websocket message.
			// n := len(c.send)
			// for i := 0; i < n; i++ {
			// 	w.Write(newline)
			// 	msgq := <-c.send
			// 	w.Write(msgq.Bytes())
			// }

			// if err := w.Close(); err != nil {
			// 	return
			// }
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) Close() error {
	close(c.send)
	return nil
}

func (c *Client) Start() error {
	go c.readPump()
	go c.writePump()
	return nil
}

func (c *Client) User() *model.User {
	return c.user
}

func (c *Client) IsSubscribed(topic string) bool {
	_, has := c.subs[topic]
	return has
}

func (c *Client) Subscribe(topic string) {
	c.subs[topic] = true
}

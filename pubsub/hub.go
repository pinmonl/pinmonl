package pubsub

import (
	"net/http"
	"time"

	"github.com/pinmonl/pinmonl/pkgs/request"
	"github.com/pinmonl/pinmonl/store"
	"github.com/sirupsen/logrus"
)

type Hub struct {
	clients   map[*Client]bool
	broadcast chan Message

	TokenSecret []byte
	TokenExpire time.Duration
	TokenIssuer string
	Users       *store.Users
}

func NewHub(tokenSecret []byte, tokenExpire time.Duration, tokenIssuer string, users *store.Users) *Hub {
	return &Hub{
		clients:   make(map[*Client]bool),
		broadcast: make(chan Message),

		TokenSecret: tokenSecret,
		TokenExpire: tokenExpire,
		TokenIssuer: tokenIssuer,
		Users:       users,
	}
}

func (h *Hub) Start() error {
	for {
		select {
		case msg := <-h.broadcast:
			for c := range h.clients {
				if !c.IsSubscribed(msg.Topic()) {
					continue
				}
				if !msg.ShouldSendTo(c) {
					continue
				}

				select {
				case c.send <- msg:
				default:
					h.Unregister(c)
				}
			}
		}
	}
	return nil
}

func (h *Hub) Broadcast(msg Message) error {
	h.broadcast <- msg
	return nil
}

func (h *Hub) Register(c *Client) error {
	h.clients[c] = true
	return nil
}

func (h *Hub) Unregister(c *Client) error {
	if _, ok := h.clients[c]; ok {
		logrus.Debugf("pubsub: hub client unregistered")
		delete(h.clients, c)
		return c.Close()
	}
	return nil
}

func (h *Hub) ServeWs() http.Handler {
	auth := request.Authenticate(h.TokenSecret, h.Users)

	fn := func(w http.ResponseWriter, r *http.Request) {
		conn, err := defaultUpgrader.Upgrade(w, r, nil)
		if err != nil {
			logrus.Debugf("pubsub: hub upgrade err(%v)", err)
			return
		}

		user := request.AuthedFrom(r.Context())
		c := NewClient(h, conn, user)
		if err := h.Register(c); err != nil {
			logrus.Debugf("pubsub: hub register err(%v)", err)
			return
		}
		if user != nil {
			logrus.Debugf("pubsub: hub client registered with user %s", user.ID)
		} else {
			logrus.Debugf("pubsub: hub client registered")
		}

		go c.Start()
	}

	return auth(http.HandlerFunc(fn))
}

var _ Pubsuber = &Hub{}

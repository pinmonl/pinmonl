package pubsub

import "github.com/pinmonl/pinmonl/logx"

// Hub maanges connected clients and message broadcasting.
type Hub struct {
	clients    map[string]map[*Client]bool
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
	quit       chan struct{}
}

// NewHub creates Hub instance.
func NewHub() *Hub {
	h := &Hub{
		clients:    make(map[string]map[*Client]bool),
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		quit:       make(chan struct{}),
	}
	go h.run()
	return h
}

func (h *Hub) run() error {
	for {
		select {
		case c := <-h.register:
			if _, ok := h.clients[c.userID]; !ok {
				h.clients[c.userID] = make(map[*Client]bool)
			}
			h.clients[c.userID][c] = true
			logx.Debugf("pubsub: user %q joined, %d clients\n", c.userID, len(h.clients[c.userID]))
		case c := <-h.unregister:
			if uc := h.clients[c.userID]; uc != nil {
				if _, ok := uc[c]; ok {
					delete(h.clients[c.userID], c)
					c.Close()
					logx.Debugf("pubsub: user %q removed\n", c.userID)
				}
			}
		case m := <-h.broadcast:
			logx.Debugf("pubsub: message broadcasting to %d", len(h.clients[m.UserID]))
			for c := range h.clients[m.UserID] {
				select {
				case c.send <- m:
				default:
				}
			}
		case <-h.quit:
			for _, uc := range h.clients {
				for c := range uc {
					c.Close()
				}
			}
			return nil
		}
	}
	return nil
}

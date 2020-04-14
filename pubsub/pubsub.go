package pubsub

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pinmonl/pinmonl/logx"
	"github.com/pinmonl/pinmonl/session"
)

// ServerOpts defines the options of creating server.
type ServerOpts struct {
	Cookie session.Store
}

// Server manages Pub/Sub service.
type Server struct {
	hub      *Hub
	upgrader websocket.Upgrader
	cookie   session.Store
}

// NewServer creates Pubsub Server instance.
func NewServer(opts *ServerOpts) *Server {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(*http.Request) bool { return true },
	}
	return &Server{
		hub:      NewHub(),
		upgrader: upgrader,
		cookie:   opts.Cookie,
	}
}

// Publish sends message to client.
func (s *Server) Publish(msg *Message) error {
	s.hub.broadcast <- msg
	return nil
}

// Handler handles web socket request.
func (s *Server) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := s.upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		sv, err := s.cookie.Get(r)
		if err != nil {
			logx.Debugf("pubsub: get cookie err(%v)", err)
			return
		}

		client := NewClient(s.hub, conn, sv.UserID)
		s.hub.register <- client
	}
}

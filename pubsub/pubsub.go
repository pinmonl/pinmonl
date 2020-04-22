package pubsub

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pinmonl/pinmonl/logx"
	"github.com/pinmonl/pinmonl/session"
	"github.com/pinmonl/pinmonl/store"
)

// ServerOpts defines the options of creating server.
type ServerOpts struct {
	SingleUser bool

	Cookie session.Store
	Users  store.UserStore
}

// Server manages Pub/Sub service.
type Server struct {
	singleUser bool

	hub      *Hub
	upgrader websocket.Upgrader
	cookie   session.Store
	users    store.UserStore
}

// NewServer creates Pubsub Server instance.
func NewServer(opts *ServerOpts) (*Server, error) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(*http.Request) bool { return true },
	}
	return &Server{
		singleUser: opts.SingleUser,

		hub:      NewHub(),
		upgrader: upgrader,
		cookie:   opts.Cookie,
		users:    opts.Users,
	}, nil
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

		userID := ""
		if s.singleUser {
			users, err := s.users.List(r.Context(), &store.UserOpts{
				ListOpts: store.ListOpts{Limit: 1},
			})
			if err != nil {
				logx.Debugf("pubsub: single user err(%v)", err)
				return
			}
			if len(users) == 0 {
				logx.Debugf("pubsub: no user found in single user mode")
				return
			}
			userID = users[0].ID
		} else {
			sv, err := s.cookie.Get(r)
			if err != nil {
				logx.Debugf("pubsub: get cookie err(%v)", err)
				return
			}
			userID = sv.UserID
		}

		client := NewClient(s.hub, conn, userID)
		s.hub.register <- client
	}
}

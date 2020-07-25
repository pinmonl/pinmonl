package pubsub

import (
	"net/http"
)

type Pubsuber interface {
	Start() error
	Register(*Client) error
	Unregister(*Client) error
	Broadcast(Message) error
	ServeWs() http.Handler
}

type Message interface {
	Topic() string
	Data() interface{}
	ShouldSendTo(*Client) bool
}

package pinl

import "github.com/pinmonl/pinmonl/pubsub"

// NewCreateMessage creates Message when created.
func NewCreateMessage(pinl Body) *pubsub.Message {
	return newPubsubMessage(pinl, "pinl.create")
}

// NewUpdateMessage creates Message when updated.
func NewUpdateMessage(pinl Body) *pubsub.Message {
	return newPubsubMessage(pinl, "pinl.update")
}

// NewDeleteMessage creates Message when deleted.
func NewDeleteMessage(pinl Body) *pubsub.Message {
	return newPubsubMessage(pinl, "pinl.delete")
}

func newPubsubMessage(pinl Body, topic string) *pubsub.Message {
	return pubsub.NewMessage(pinl.UserID, topic, pinl)
}

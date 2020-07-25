package message

import (
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/pubsub"
)

type PinlUpdated struct {
	pinl *model.Pinl
}

func NewPinlUpdated(pinl *model.Pinl) *PinlUpdated {
	return &PinlUpdated{pinl: pinl}
}

func (p *PinlUpdated) Topic() string { return "pinl_updated" }

func (p *PinlUpdated) Data() interface{} { return p.pinl }

func (p *PinlUpdated) ShouldSendTo(c *pubsub.Client) bool {
	if c.User() == nil {
		return false
	}
	return c.User().ID == p.pinl.UserID
}

var _ pubsub.Message = &PinlUpdated{}

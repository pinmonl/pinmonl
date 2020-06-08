package model

type Sharepin struct {
	ID      string      `json:"id"`
	ShareID string      `json:"shareId"`
	PinlID  string      `json:"pinlId"`
	Status  ShareStatus `json:"status"`

	Share *Share `json:"share,omitempty"`
	Pinl  *Pinl  `json:"pinl,omitempty"`
}

type SharepinList []*Sharepin

func (sl SharepinList) PinlKeys() []string {
	keys := make([]string, len(sl))
	for i, s := range sl {
		keys[i] = s.PinlID
	}
	return keys
}

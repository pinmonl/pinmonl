package model

type Sharepin struct {
	ID      string `json:"id"`
	ShareID string `json:"shareId"`
	PinlID  string `json:"pinlId"`
	Status  Status `json:"status"`

	Share *Share `json:"share,omitempty"`
	Pinl  *Pinl  `json:"pinl,omitempty"`
}

type SharepinList []*Sharepin

func (sl SharepinList) Pinls() PinlList {
	pinls := make([]*Pinl, len(sl))
	for i := range sl {
		pinls[i] = sl[i].Pinl
	}
	return pinls
}

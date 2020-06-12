package model

type Pinpkg struct {
	ID     string `json:"id"`
	PinlID string `json:"pinlId"`
	PkgID  string `json:"pkgId"`

	Pinl *Pinl `json:"pinl,omitempty"`
	Pkg  *Pkg  `json:"pkg,omitempty"`
}

type PinpkgList []*Pinpkg

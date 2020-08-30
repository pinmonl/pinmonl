package model

type Pinpkg struct {
	ID     string `json:"id"`
	PinlID string `json:"pinlId"`
	PkgID  string `json:"pkgId"`

	Pinl *Pinl `json:"pinl,omitempty"`
	Pkg  *Pkg  `json:"pkg,omitempty"`
}

type PinpkgList []*Pinpkg

func (pl PinpkgList) Pkgs() PkgList {
	out := make([]*Pkg, len(pl))
	for i := range pl {
		out[i] = pl[i].Pkg
	}
	return out
}

func (pl PinpkgList) PkgKeys() []string {
	out := make([]string, len(pl))
	for i := range pl {
		out[i] = pl[i].PkgID
	}
	return out
}

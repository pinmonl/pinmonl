package model

type Monpkg struct {
	ID     string     `json:"id"`
	MonlID string     `json:"monlId"`
	PkgID  string     `json:"pkgId"`
	Kind   MonpkgKind `json:"kind"`

	Monl *Monl `json:"monl,omitempty"`
	Pkg  *Pkg  `json:"pkg,omitempty"`
}

type MonpkgKind int

type MonpkgList []*Monpkg

func (ml MonpkgList) Pkgs() PkgList {
	out := make([]*Pkg, len(ml))
	for i := range ml {
		out[i] = ml[i].Pkg
	}
	return out
}

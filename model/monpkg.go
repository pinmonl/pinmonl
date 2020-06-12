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

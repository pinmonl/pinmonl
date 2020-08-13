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

const (
	MonpkgDirect MonpkgKind = iota
	MonpkgDerived
)

type MonpkgList []*Monpkg

func (ml MonpkgList) Pkgs() PkgList {
	out := make([]*Pkg, len(ml))
	for i := range ml {
		out[i] = ml[i].Pkg
	}
	return out
}

func (ml MonpkgList) PkgsByMonl() map[string]PkgList {
	out := make(map[string]PkgList)
	for i := range ml {
		k := ml[i].MonlID
		out[k] = append(out[k], ml[i].Pkg)
	}
	return out
}

package model

import "github.com/pinmonl/pinmonl/model/field"

type Pinl struct {
	ID          string     `json:"id"`
	UserID      string     `json:"userId"`
	MonlID      string     `json:"monlId"`
	URL         string     `json:"url"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	ImageID     string     `json:"imageId"`
	Status      Status     `json:"status"`
	HasPinpkgs  bool       `json:"hasPinpkgs"`
	CreatedAt   field.Time `json:"createdAt"`
	UpdatedAt   field.Time `json:"updatedAt"`

	Tags     *TagList  `json:"tags,omitempty"`
	TagNames *[]string `json:"tagNames,omitempty"`
	Pkgs     *PkgList  `json:"pkgs,omitempty"`
	PkgIDs   *[]string `json:"pkgIds,omitempty"`
}

func (p Pinl) MorphKey() string  { return p.ID }
func (p Pinl) MorphName() string { return "pinl" }

func (p *Pinl) SetTagNames(tags TagList) {
	tn := tags.Names()
	p.TagNames = &tn
}

func (p *Pinl) SetPkgs(pkgs PkgList) {
	p.Pkgs = &pkgs
}

func (p *Pinl) SetPkgIDs(pkgs PkgList) {
	ids := pkgs.Keys()
	p.PkgIDs = &ids
}

type PinlList []*Pinl

func (pl PinlList) Keys() []string {
	keys := make([]string, len(pl))
	for i := range pl {
		keys[i] = pl[i].ID
	}
	return keys
}

func (pl PinlList) Morphables() MorphableList {
	list := make([]Morphable, len(pl))
	for i := range pl {
		list[i] = pl[i]
	}
	return list
}

func (pl PinlList) SetTagNames(tagMap map[string]TagList) {
	for i := range pl {
		k := pl[i].ID
		pl[i].SetTagNames(tagMap[k])
	}
}

func (pl PinlList) MonlKeys() []string {
	keys := make([]string, len(pl))
	for i := range pl {
		keys[i] = pl[i].MonlID
	}
	return keys
}

func (pl PinlList) SetPkgs(pMap map[string]PkgList) {
	for i := range pl {
		k := pl[i].MonlID
		pkgs := append(make([]*Pkg, 0), pMap[k]...)
		pl[i].SetPkgs(pkgs)
	}
}

func (pl PinlList) SetPkgIDs(pMap map[string]PkgList) {
	for i := range pl {
		k := pl[i].ID
		pl[i].SetPkgIDs(pMap[k])
	}
}

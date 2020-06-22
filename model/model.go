package model

type Morphable interface {
	MorphKey() string
	MorphName() string
}

type MorphableList []Morphable

func (ml MorphableList) IsMixed() bool {
	checks := make(map[string]int)
	for _, m := range ml {
		checks[m.MorphName()]++
	}
	return len(checks) > 1
}

func (ml MorphableList) MorphName() string {
	if len(ml) == 0 {
		return ""
	}
	return ml[0].MorphName()
}

func (ml MorphableList) MorphKeys() []string {
	keys := make([]string, len(ml))
	for i := range ml {
		keys[i] = ml[i].MorphKey()
	}
	return keys
}

type Status int

const (
	Active Status = iota
	Preparing
	Deleted
)

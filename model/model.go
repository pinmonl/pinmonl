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
	var keys []string
	for _, m := range ml {
		keys = append(keys, m.MorphKey())
	}
	return keys
}

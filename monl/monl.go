package monl

import "time"

type Monler interface {
	MonlerName() string
	Open(url *URL) (Object, error)
}

type Object interface {
	ObjectName() string
	Stats() []StatGroup
	Aliases() []string
	Related() []string
}

type StatGroup interface {
	GroupName() string
	IsAppend() bool
	Stats() []Stat
}

type Stat interface {
	At() time.Time
	Value() string
	Substats() []Stat
}

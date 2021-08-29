package monl

import (
	"time"

	"github.com/pinmonl/pinmonl/monlurl"
)

type Monler interface {
	MonlerName() string
	Open(url *monlurl.Url) (Object, error)
}

type RemoteResource interface {
	Populate() error
	Close() error
}

type Object interface {
	RemoteResource
	ObjectName() string
	Stats() []StatGroup
}

type StatGroup interface {
	RemoteResource
	GroupName() string
	IsAppend() bool
	Stats() []Stat
	Latest() []Stat
}

type Stat interface {
	At() time.Time
	Value() string
	Substats() []StatGroup
	IsLatest() bool
	Hash() string
}

type SimpleObject struct {
	objectName string
	stats      []StatGroup
	aliases    []string
	related    []string
}

func (s SimpleObject) Populate() error    { return nil }
func (s SimpleObject) Close() error       { return nil }
func (s SimpleObject) ObjectName() string { return s.objectName }
func (s SimpleObject) Stats() []StatGroup { return s.stats }
func (s SimpleObject) Aliases() []string  { return s.aliases }
func (s SimpleObject) Related() []string  { return s.related }

type SimpleStat struct {
	at       time.Time
	value    string
	substats []StatGroup
	isLatest bool
	hash     string
}

func (s SimpleStat) At() time.Time         { return s.at }
func (s SimpleStat) Value() string         { return s.value }
func (s SimpleStat) Substats() []StatGroup { return s.substats }
func (s SimpleStat) IsLatest() bool        { return s.isLatest }
func (s SimpleStat) Hash() string          { return s.hash }

type SimpleStatGroup struct {
	groupName string
	isAppend  bool
	stats     []Stat
	latest    []Stat
}

func (s SimpleStatGroup) Populate() error   { return nil }
func (s SimpleStatGroup) Close() error      { return nil }
func (s SimpleStatGroup) GroupName() string { return s.groupName }
func (s SimpleStatGroup) IsAppend() bool    { return s.isAppend }
func (s SimpleStatGroup) Stats() []Stat     { return s.stats }
func (s SimpleStatGroup) Latest() []Stat    { return s.latest }

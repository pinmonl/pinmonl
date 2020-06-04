package model

import "github.com/pinmonl/pinmonl/model/field"

type Stat struct {
	ID          string        `json:"id"`
	PkgID       string        `json:"pkgId"`
	ParentID    string        `json:"parentId"`
	RecordedAt  field.Time    `json:"recordedAt"`
	Kind        StatKind      `json:"kind"`
	Value       string        `json:"value"`
	ValueType   StatValueType `json:"valueType"`
	Checksum    string        `json:"digest"`
	Weight      int           `json:"weight"`
	IsLatest    bool          `json:"isLatest"`
	HasChildren bool          `json:"hasChildren"`

	Substats *[]*Stat `json:"substats,omitempty"`
}

func (s Stat) MorphKey() string  { return s.ID }
func (s Stat) MorphName() string { return "stat" }

type StatKind int

const (
	StatKindNotSet StatKind = iota
	TagStat
	AliasStat
	StarStat
	ForkStat
	OpenIssueStat
	LangStat
	FileCountStat
	DownloadStat
	PullStat
)

type StatValueType int

const (
	StatValueTypeNotSet StatValueType = iota
	StringStat
	IntegerStat
)

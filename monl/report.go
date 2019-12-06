package monl

type Report interface {
	RawUrl() string
	Uri() string
	Vendor() Vendor
	Popularity() StatCollection
	Latest() Stat
	Next() Stat
	NextAll() StatCollection
	Length() int
}

type SimpleReport struct {
	rawUrl string
	vendor Vendor
	stats  []Stat
}

func NewReport(
	rawUrl string,
	vendor Vendor,
	stats []Stat,
) *SimpleReport {
	return &SimpleReport{
		rawUrl: rawUrl,
		vendor: vendor,
		stats:  stats,
	}
}

func (s *SimpleReport) RawUrl() string { return s.rawUrl }

func (s *SimpleReport) Vendor() Vendor { return s.vendor }

func (s *SimpleReport) Stats() []Stat { return s.stats }

func (s *SimpleReport) Length() int { return len(s.stats) }

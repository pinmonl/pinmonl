package model

// Monpkg stores the relation between Monl and Pkg.
type Monpkg struct {
	*Monl
	*Pkg
	MonlID string    `json:"monlId" db:"monpkg_monl_id"`
	PkgID  string    `json:"pkgId"  db:"monpkg_pkg_id"`
	Tie    MonpkgTie `json:"tie"    db:"monpkg_tie"`
}

// MonpkgTie defines the relation tie.
type MonpkgTie int

const (
	// MonpkgTiePrimary indicates the primary relation tie.
	MonpkgTiePrimary MonpkgTie = iota
	// MonpkgTieSecondary indicates the secondary relation tie.
	MonpkgTieSecondary
)

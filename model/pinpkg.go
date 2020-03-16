package model

// Pinpkg defines the connection between Pinl and Pkg.
type Pinpkg struct {
	*Pinl
	*Pkg
	PinlID string `json:"pinlId" db:"pinpkg_pinl_id"`
	PkgID  string `json:"pkgId"  db:"pinpkg_pkg_id"`
	Sort   int64  `json:"sort"   db:"pinpkg_sort"`
}

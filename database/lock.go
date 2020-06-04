package database

import "sync"

type Locker interface {
	sync.Locker
}

type NopLocker struct{}

func (n *NopLocker) Lock()   {}
func (n *NopLocker) Unlock() {}

func NewDriverLocker(driver string) Locker {
	switch driver {
	case "sqlite3":
		return &sync.Mutex{}
	default:
		return &NopLocker{}
	}
}

package repository

import (
	"sync"
)

var (
	internalManager *Manager
	once sync.Once
)

func InitManager(opt ...Option) *Manager {
	once.Do(func() {
		internalManager = NewManager(opt...)
	})
	return internalManager
}

package manager

import (
	"basicLoginRest/internal/session"
	"sync"
)

var (
	internalManager session.Manager
	once sync.Once
)

func InitManager(opt ...Option) session.Manager {
	once.Do(func() {
		internalManager = newManager(opt...)
	})
	return internalManager
}

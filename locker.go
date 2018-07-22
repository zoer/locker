package locker

import (
	"context"
)

// Locker represents interface to create a lock
type Locker interface {
	Lock(context.Context, ...Option) (Lock, error)
}

// Lock represents interface to unlock existed lock
type Lock interface {
	Unlock()
}

// Option represents the lock option
type Option func(*options)

package locker

import "time"

type options struct {
	ttl     time.Duration
	waitTTL time.Duration
	key     string
}

func newOptions() *options {
	return &options{}
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

func (o options) withOptions(opts ...Option) options {
	o.apply(opts...)
	return o
}

// WithTTL sets the TTL option.
func WithTTL(ttl time.Duration) Option {
	return func(o *options) {
		if ttl > 0 {
			o.ttl = ttl
		}
	}
}

// WithWaitTTL sets the wait for receiving the lock TTL.
func WithWaitTTL(ttl time.Duration) Option {
	return func(o *options) {
		if ttl > 0 {
			o.waitTTL = ttl
		}
	}
}

// WithKey sets the lock's key
func WithKey(key string) Option {
	return func(o *options) {
		o.key = key
	}
}

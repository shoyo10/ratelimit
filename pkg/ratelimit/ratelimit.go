package ratelimit

import (
	"context"
	"time"
)

type Config struct {
	// a time interval to allow max resource
	Interval time.Duration `yaml:"interval" mapstructure:"interval"`

	// max request within a time interval
	MaxRequest int64 `yaml:"max_request" mapstructure:"max_request"`
}

type RequestInfo struct {
	Count       int64
	IsAvaliable bool
}

type Limiter interface {
	// Avaliable it will try to take a resource and return true if it is enough, or return false.
	Avaliable(ctx context.Context, target string) (RequestInfo, error)
}

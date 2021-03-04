package redis

import (
	"context"
	"ratelimit/pkg/ratelimit"
	"ratelimit/pkg/redis"
	"ratelimit/pkg/zerolog"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
)

func TestIsAvaliable(t *testing.T) {
	zerolog.Init(&zerolog.Config{
		Debug: true,
		Local: true,
		AppID: "test",
		Env:   "test",
	})
	logger := log.With().Logger()
	ctx := logger.WithContext(context.Background())

	rdb := redis.New(redis.Config{
		Addr: "127.0.0.1:6379",
	})
	cfg := ratelimit.Config{
		Interval:   2 * time.Second,
		MaxRequest: 5,
	}
	limiter := NewRateLimit(cfg, rdb)

	target := "127.0.0.1"

	testCases := []struct {
		elapsedTime time.Duration
		want        ratelimit.RequestInfo
	}{
		{
			elapsedTime: 1 * time.Millisecond,
			want: ratelimit.RequestInfo{
				Count:       1,
				IsAvaliable: true,
			},
		},
		{
			elapsedTime: 1 * time.Millisecond,
			want: ratelimit.RequestInfo{
				Count:       2,
				IsAvaliable: true,
			},
		},
		{
			elapsedTime: 1 * time.Millisecond,
			want: ratelimit.RequestInfo{
				Count:       3,
				IsAvaliable: true,
			},
		},
		{
			elapsedTime: 500 * time.Millisecond,
			want: ratelimit.RequestInfo{
				Count:       4,
				IsAvaliable: true,
			},
		},
		{
			elapsedTime: 1 * time.Millisecond,
			want: ratelimit.RequestInfo{
				Count:       5,
				IsAvaliable: true,
			},
		},
		{
			elapsedTime: 1 * time.Millisecond,
			want: ratelimit.RequestInfo{
				Count:       5,
				IsAvaliable: false,
			},
		},
		{
			elapsedTime: 500 * time.Millisecond,
			want: ratelimit.RequestInfo{
				Count:       5,
				IsAvaliable: false,
			},
		},
		{
			elapsedTime: 2 * time.Second,
			want: ratelimit.RequestInfo{
				Count:       1,
				IsAvaliable: true,
			},
		},
		{
			elapsedTime: 3 * time.Second,
			want: ratelimit.RequestInfo{
				Count:       1,
				IsAvaliable: true,
			},
		},
		{
			elapsedTime: 300 * time.Millisecond,
			want: ratelimit.RequestInfo{
				Count:       2,
				IsAvaliable: true,
			},
		},
	}

	for _, tc := range testCases {
		time.Sleep(tc.elapsedTime)
		got, err := limiter.Avaliable(ctx, target)
		if err != nil {
			t.Errorf("got error: %v", err)
			continue
		}
		if got.IsAvaliable != tc.want.IsAvaliable || got.Count != tc.want.Count {
			t.Errorf("Avaliable() = %v, want: %v", got, tc.want)
		}
	}
}

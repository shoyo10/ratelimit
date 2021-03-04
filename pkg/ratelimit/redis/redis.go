package redis

import (
	"context"
	"ratelimit/pkg/ratelimit"
	"ratelimit/pkg/tokenbucket"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

const maxRetries = 1000

type limiter struct {
	cfg ratelimit.Config
	rdb *redis.Client
}

type request struct {
	Bucket    *tokenbucket.Bucket
	ExpiredAt time.Time
	Count     int64
}

// NewRateLimit create redis rate limit instance
func NewRateLimit(cfg ratelimit.Config, rdb *redis.Client) ratelimit.Limiter {
	return &limiter{
		cfg: cfg,
		rdb: rdb,
	}
}

func (l *limiter) Avaliable(ctx context.Context, target string) (ratelimit.RequestInfo, error) {
	var resp ratelimit.RequestInfo

	// redis Transactional function.
	txf := func(tx *redis.Tx) error {
		data, err := tx.Get(ctx, target).Result()
		if err != nil && err != redis.Nil {
			return err
		}

		req := &request{}
		if err == redis.Nil {
			req = l.newRequest()
		} else {
			err = json.Unmarshal([]byte(data), req)
			if err != nil {
				return err
			}
		}

		if time.Now().UTC().After(req.ExpiredAt) {
			req = l.newRequest()
		}

		if req.Count == l.cfg.MaxRequest {
			resp.IsAvaliable = false
			resp.Count = req.Count
		} else {
			resp.IsAvaliable = req.Bucket.TokenAvaliable()
			if resp.IsAvaliable {
				req.Count++
			}
			resp.Count = req.Count
		}

		b, err := json.Marshal(req)
		if err != nil {
			return err
		}

		// Operation is commited only if the watched keys remain unchanged.
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			return pipe.Set(ctx, target, string(b), l.cfg.Interval).Err()
		})
		return err
	}

	for i := 0; i < maxRetries; i++ {
		err := l.rdb.Watch(ctx, txf, target)
		if err == nil {
			return resp, nil
		}
		if err == redis.TxFailedErr {
			// Optimistic lock lost. Retry.
			log.Ctx(ctx).Info().Msgf("%v", err)
			continue
		}

		return resp, err
	}

	return resp, errors.New("reached maximum number of retries")
}

func (l *limiter) newRequest() *request {
	return &request{
		Bucket:    l.newTokenBucket(),
		ExpiredAt: time.Now().UTC().Add(l.cfg.Interval),
	}
}

func (l *limiter) newTokenBucket() *tokenbucket.Bucket {
	return tokenbucket.New(tokenbucket.Config{
		Rate:         l.cfg.Interval / time.Duration(l.cfg.MaxRequest),
		MaxBurstSize: l.cfg.MaxRequest,
	})
}

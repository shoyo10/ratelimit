package tokenbucket

import (
	"sync"
	"time"
)

type Config struct {
	// a time interval to generate a new token into the bucket
	Rate time.Duration

	// max number of tokens can be consumed within a time interval
	MaxBurstSize int64
}

type Bucket struct {
	sync.Mutex
	Cfg           Config
	CurrentTokens int64
	LatestAddedAt time.Time
}

// New create token bucket, initial bucket tokens are full
func New(cfg Config) *Bucket {
	return &Bucket{
		Cfg:           cfg,
		CurrentTokens: cfg.MaxBurstSize,
		LatestAddedAt: time.Now().UTC(),
	}
}

// TokenAvaliable if there are enough tokens, it remove one token and return true, otherwise return false
func (b *Bucket) TokenAvaliable() bool {
	b.addTokens()
	return b.removeToken()
}

func (b *Bucket) addTokens() {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()

	nowTime := time.Now().UTC()
	elapsedTime := nowTime.Sub(b.LatestAddedAt)
	newTokens := int64(elapsedTime / b.Cfg.Rate)
	if newTokens == 0 {
		return
	}

	b.CurrentTokens += newTokens
	if b.CurrentTokens > b.Cfg.MaxBurstSize {
		b.CurrentTokens = b.Cfg.MaxBurstSize
	}

	b.LatestAddedAt = b.LatestAddedAt.Add(time.Duration(newTokens) * b.Cfg.Rate)
}

func (b *Bucket) removeToken() bool {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()

	if b.CurrentTokens == 0 {
		return false
	}
	b.CurrentTokens--
	return true
}

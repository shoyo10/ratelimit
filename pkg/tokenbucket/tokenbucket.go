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

// TokenAvaliable remove one token if there are enough tokens.
func (b *Bucket) TokenAvaliable() bool {
	return b.TokenAvaliableN(1)
}

// TokenAvaliableN if there are enough tokens, it will remove tokens and return true, otherwise return false
func (b *Bucket) TokenAvaliableN(amount int64) bool {
	b.addTokens()
	return b.removeTokenN(amount)
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

func (b *Bucket) removeTokenN(amount int64) bool {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()

	if b.CurrentTokens == 0 || amount > b.CurrentTokens {
		return false
	}
	b.CurrentTokens -= amount
	return true
}

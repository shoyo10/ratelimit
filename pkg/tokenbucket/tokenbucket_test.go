package tokenbucket

import (
	"testing"
	"time"
)

func TestTokenAvaliable(t *testing.T) {
	cfg := Config{
		Rate:         500 * time.Millisecond,
		MaxBurstSize: 2,
	}
	tb := New(cfg)
	testCases := []struct {
		elapsedTime time.Duration
		want        bool
	}{
		{
			elapsedTime: 1 * time.Millisecond,
			want:        true,
		},
		{
			elapsedTime: 1 * time.Millisecond,
			want:        true,
		},
		{
			elapsedTime: 1 * time.Millisecond,
			want:        false,
		},
		{
			elapsedTime: 500 * time.Millisecond,
			want:        true,
		},
		{
			elapsedTime: 300 * time.Millisecond,
			want:        false,
		},
		{
			elapsedTime: 200 * time.Millisecond,
			want:        true,
		},
		{
			elapsedTime: 2000 * time.Millisecond,
			want:        true,
		},
	}

	for _, tc := range testCases {
		time.Sleep(tc.elapsedTime)
		if got := tb.TokenAvaliable(); got != tc.want {
			t.Errorf("TokenAvaliable() = %v, want: %v", got, tc.want)
		}
	}
}

func TestTokenAvaliableN(t *testing.T) {
	cfg := Config{
		Rate:         500 * time.Millisecond,
		MaxBurstSize: 4,
	}
	tb := New(cfg)
	testCases := []struct {
		elapsedTime time.Duration
		amount      int64
		want        bool
	}{
		{
			elapsedTime: 1 * time.Millisecond,
			amount:      2,
			want:        true,
		},
		{
			elapsedTime: 1 * time.Millisecond,
			amount:      3,
			want:        false,
		},
		{
			elapsedTime: 1 * time.Millisecond,
			amount:      2,
			want:        true,
		},
		{
			elapsedTime: 500 * time.Millisecond,
			amount:      1,
			want:        true,
		},
		{
			elapsedTime: 300 * time.Millisecond,
			amount:      1,
			want:        false,
		},
		{
			elapsedTime: 200 * time.Millisecond,
			amount:      1,
			want:        true,
		},
	}

	for _, tc := range testCases {
		time.Sleep(tc.elapsedTime)
		if got := tb.TokenAvaliableN(tc.amount); got != tc.want {
			t.Errorf("TokenAvaliableN() = %v, want: %v", got, tc.want)
		}
	}
}

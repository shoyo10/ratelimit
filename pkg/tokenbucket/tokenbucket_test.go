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

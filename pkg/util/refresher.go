package util

import (
	"time"
)

type Refresher struct {
	interval	time.Duration
	ticker		*time.Ticker
	stopChan	chan struct{}
	onRefresh	func()
} 

func RefresherCreate(interval time.Duration, onRefresh func()) *Refresher {
	return &Refresher{
		interval: interval,
		stopChan: make(chan struct{}),
		onRefresh: onRefresh,
	}
}

func (r *Refresher) Start() {
	r.ticker = time.NewTicker(r.interval)
	go func() {
		for {
			select {
			case <-r.ticker.C:
				r.onRefresh()
			case <-r.stopChan:
				r.ticker.Stop()
				return
			}
		}
	}()
}

func (r *Refresher) Stop() {
	close(r.stopChan)
}
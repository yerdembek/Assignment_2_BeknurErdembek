package worker

import (
	"log"
	"time"
)

type StatsProvider interface {
	Requests() int64
	KeysCount() int
}

type Worker struct {
	stats  StatsProvider
	ticker *time.Ticker
	stopCh chan struct{}
	doneCh chan struct{}
}

func New(stats StatsProvider) *Worker {
	return &Worker{
		stats:  stats,
		stopCh: make(chan struct{}),
		doneCh: make(chan struct{}),
	}
}

func (w *Worker) Start() {
	w.ticker = time.NewTicker(5 * time.Second)
	go func() {
		defer close(w.doneCh)
		for {
			select {
			case <-w.ticker.C:
				reqs := w.stats.Requests()
				keys := w.stats.KeysCount()
				log.Printf("[worker] requests=%d keys=%d\n", reqs, keys)
			case <-w.stopCh:
				w.ticker.Stop()
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	close(w.stopCh)
	<-w.doneCh
}

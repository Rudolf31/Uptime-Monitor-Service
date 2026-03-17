package worker

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"
	"uptime-monitor/internal/storage"
)

type Pool struct {
	jobs    chan int
	wg      sync.WaitGroup
	storage storage.Storage
	client  *http.Client
}

func NewPool(storage storage.Storage) *Pool {
	return &Pool{
		jobs:    make(chan int),
		storage: storage,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (p *Pool) Start(ctx context.Context) {
	for i := 0; i < 5; i++ {
		p.wg.Add(1)
		go p.worker(ctx)
	}
}

func (p *Pool) worker(ctx context.Context) {
	defer p.wg.Done()

	for {
		select {
		case id, ok := <-p.jobs:
			if !ok {
				return
			}
			p.check(id)

		case <-ctx.Done():
			return
		}
	}
}

func (p *Pool) check(id int) {
	monitor, err := p.storage.GetByID(id)
	if err != nil {
		slog.Error("monitor not found", "error", err)
		return
	}

	status, duration := CheckURL(monitor.URL, p.client)

	err = p.storage.UpdateCheckResult(
		id,
		status,
		time.Now(),
		duration.Milliseconds(),
	)
	if err != nil {
		slog.Error("failed to update monitor", "error", err)
	}
}

func (p *Pool) Jobs() chan int {
	return p.jobs
}

func (p *Pool) Stop() {
	close(p.jobs)
	p.wg.Wait()
}

package worker

import (
	"context"
	"log/slog"
	"time"
	"uptime-monitor/internal/storage"
)

type Scheduler struct {
	storage storage.Storage
	jobs    chan<- int
}

func NewScheduler(storage storage.Storage, jobs chan<- int) *Scheduler {
	return &Scheduler{
		storage: storage,
		jobs:    jobs, // Need pool.jobs
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ticker.C:
			s.CheckTime()
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (s *Scheduler) CheckTime() {
	now := time.Now()
	monitors, err := s.storage.GetAll()
	if err != nil {
		slog.Error("failed to get monitors", "error", err)
	}

	for _, monitor := range monitors {
		if now.After(monitor.NextCheck) {
			nextCheck := now.Add(time.Duration(monitor.Interval) * time.Second)
			s.storage.UpdateNextCheck(monitor.ID, nextCheck)

			s.jobs <- monitor.ID
		}
	}
}

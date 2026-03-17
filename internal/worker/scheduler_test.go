package worker

import (
	"testing"
	"time"
	"uptime-monitor/internal/models"
	"uptime-monitor/internal/storage"
)

func TestSchedulerCheckTime(t *testing.T) {
	store := storage.NewMemoryStorage()
	jobs := make(chan int, 1)

	scheduler := NewScheduler(store, jobs)

	monitor := &models.Monitor{
		URL:       "https://example.com",
		Interval:  10,
		NextCheck: time.Now().Add(-1 * time.Second),
	}

	store.Create(monitor)

	scheduler.CheckTime()

	var id int
	select {
	case id = <-jobs:
		if id != monitor.ID {
			t.Fatalf("expected id %d, got %d", monitor.ID, id)
		}
	default:
		t.Fatal("expected job from scheduler")
	}

	newMonitor, err := store.GetByID(id)
	if err != nil {
		t.Fatal(err)
	}
	if newMonitor.NextCheck.Before(time.Now()) {
		t.Error("expected nextCheck time after now")
	}
}

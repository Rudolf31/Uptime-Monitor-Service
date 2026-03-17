package storage

import (
	"fmt"
	"sync"
	"testing"
	"uptime-monitor/internal/models"
)

func TestStorageConcurrentCreate(t *testing.T) {
	store := NewMemoryStorage()

	const count = 1000
	var wg sync.WaitGroup

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			m := &models.Monitor{
				URL:      fmt.Sprintf("https://example.com/%d", i),
				Interval: 10,
			}
			store.Create(m)
		}(i)
	}

	wg.Wait()

	monitors, _ := store.GetAll()
	if len(monitors) != count {
		t.Fatalf("expected %d monitors, got %d", count, len(monitors))
	}
}

func TestStorageConcurrentMixed(t *testing.T) {
	store := NewMemoryStorage()

	const count = 1000
	var wg sync.WaitGroup

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			m := &models.Monitor{
				URL:      fmt.Sprintf("https://example.com/%d", i),
				Interval: 10,
			}
			store.Create(m)
		}(i)
	}

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			store.GetAll()
		}()
	}

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			store.Delete(i)
		}()
	}

	wg.Wait()
}

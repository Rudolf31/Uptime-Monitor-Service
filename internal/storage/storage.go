package storage

import (
	"sync"
	"uptime-monitor/internal/models"
)

type MemoryStorage struct {
	mu       sync.RWMutex
	monitors map[int]*models.Monitor
	nextID   int
}

func NewStorage() *MemoryStorage {
	return &MemoryStorage{
		monitors: make(map[int]*models.Monitor),
		nextID:   1,
	}
}

func (s *MemoryStorage) Create(monitor *models.Monitor) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	monitor.ID = s.nextID
	s.nextID++

	s.monitors[monitor.ID] = monitor
	return nil
}

func (s *MemoryStorage) GetAll() ([]*models.Monitor, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	monitors := make([]*models.Monitor, 0, len(s.monitors))
	for _, m := range s.monitors {
		monitors = append(monitors, m)
	}

	return monitors, nil
}

func (s *MemoryStorage) GetByID(id int) (*models.Monitor, error) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	monitor, exists := s.monitors[id]
	if !exists {
		return nil, ErrMonitorNotFound
	}

	return monitor, nil
}

func (s *MemoryStorage) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.monitors[id]
	if !exists {
		return ErrMonitorNotFound
	}

	return nil
}

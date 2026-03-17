package storage

import (
	"time"
	"uptime-monitor/internal/models"
)

type Storage interface {
	Create(m *models.Monitor) error
	GetAll() ([]*models.Monitor, error)
	GetByID(id int) (*models.Monitor, error)
	Delete(id int) error
	UpdateCheckResult(id int, status string, lastCheck time.Time, responseTime int64) error
	UpdateNextCheck(id int, newTime time.Time) error
}

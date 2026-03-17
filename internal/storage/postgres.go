package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"uptime-monitor/internal/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sqlx.DB
}

func NewPostgresStorage(db *sqlx.DB) *PostgresStorage {
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	return &PostgresStorage{db: db}
}

func (p *PostgresStorage) Create(m *models.Monitor) error {
	query := `
	INSERT INTO monitors (url, interval, status, next_check, response_time)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`

	return p.db.QueryRowx(
		query,
		m.URL,
		m.Interval,
		m.Status,
		m.NextCheck,
		m.ResponseTime,
	).Scan(&m.ID)
}

func (p *PostgresStorage) GetAll() ([]*models.Monitor, error) {
	query := `SELECT id, url, interval, status, last_check, next_check, response_time FROM monitors`

	var monitors []*models.Monitor
	if err := p.db.Select(&monitors, query); err != nil {
		return nil, err
	}

	return monitors, nil
}

func (p *PostgresStorage) GetByID(id int) (*models.Monitor, error) {
	query := `
		SELECT id, url, interval, status, last_check, next_check, response_time
		FROM monitors
		WHERE id = $1
		`
	var monitor models.Monitor
	if err := p.db.Get(&monitor, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMonitorNotFound
		}
		return nil, err
	}

	var noteRows []models.Note
	notesQuery := `
		SELECT status, check_time, response_time
		FROM notes
		WHERE monitor_id = $1
		ORDER BY check_time DESC
		LIMIT 3
`
	if err := p.db.Select(&noteRows, notesQuery, id); err != nil {
		return nil, err
	}

	history := make([]*models.Note, 0, len(noteRows))
	for i := range noteRows {
		n := noteRows[i]
		history = append(history, &n)
	}
	monitor.History = history

	return &monitor, nil
}

func (p *PostgresStorage) Delete(id int) error {
	query := `DELETE FROM monitors WHERE id = $1`
	res, err := p.db.Exec(query, id)
	if err != nil {
		return err
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return ErrMonitorNotFound
	}
	return nil
}

func (p *PostgresStorage) UpdateCheckResult(id int, status string, lastCheck time.Time, responseTime int64) error {
	tx, err := p.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE monitors SET status=$1, last_check = $2, response_time=$3 WHERE id=$4`
	res, err := p.db.Exec(query, status, lastCheck, responseTime, id)
	if err != nil {
		return nil
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrMonitorNotFound
	}

	query = `INSERT INTO notes (monitor_id, status, check_time, response_time) VALUES ($1, $2, $3, $4)`
	res, err = p.db.Exec(query, id, status, lastCheck, responseTime)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrMonitorNotFound
	}

	return tx.Commit()

}

func (p *PostgresStorage) UpdateNextCheck(id int, newTime time.Time) error {
	res, err := p.db.Exec(`UPDATE monitors SET next_check=$1 WHERE id=$2`, newTime, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrMonitorNotFound
	}
	return nil
}

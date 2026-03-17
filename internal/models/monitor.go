package models

import (
	"net/url"
	"time"
)

type Monitor struct {
	ID           int        `json:"id" db:"id"`
	URL          string     `json:"url" db:"url"`
	Interval     int        `json:"interval" db:"interval"`
	Status       string     `json:"status" db:"status"`
	LastCheck    *time.Time `json:"last_check,omitempty" db:"last_check"`
	NextCheck    *time.Time `json:"next_check,omitempty" db:"next_check"`
	ResponseTime *int64     `json:"response_time,omitempty" db:"response_time"`
	History      []*Note    `json:"history,omitempty" db:"-"`
}

type MonitorResponse struct {
	ID       int    `json:"id"`
	URL      string `json:"url"`
	Interval int    `json:"interval"`
	Status   string `json:"status"`
}

func (m *Monitor) Validate() error {
	if m.URL == "" {
		return ErrInvalidData
	}

	parsed, err := url.ParseRequestURI(m.URL)
	if err != nil {
		return ErrInvalidData
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return ErrInvalidData
	}

	if m.Interval < 30 {
		m.Interval = 30
	}

	return nil
}

type CreateMonitorRequest struct {
	URL      string `json:"url"`
	Interval int    `json:"interval"`
	Status   string `json:"status"`
}

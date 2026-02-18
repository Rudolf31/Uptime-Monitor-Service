package models

import "strings"

type Monitor struct {
	ID       int    `json:"id"`
	URL      string `json:"url"`
	Interval int    `json:"interval"`
	Status   string `json:"status"`
}

func (m *Monitor) Validate() error {
	if m.URL == "" {
		return ErrInvalidData
	}
	if !strings.HasPrefix(m.URL, "http://") && !strings.HasPrefix(m.URL, "https://") {
		return ErrInvalidData
	}
	if m.Interval <= 0 {
		return ErrInvalidData
	}
	return nil
}

type CreateMonitorRequest struct {
	URL      string `json:"url"`
	Interval int    `json:"interval"`
	Status   string `json:"status"`
}

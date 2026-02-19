package models

import "net/url"

type Monitor struct {
	ID       int     `json:"id"`
	URL      string  `json:"url"`
	Interval int     `json:"interval"`
	Status   string  `json:"status"`
	History  []*Note `json:"history,omitempty"`
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

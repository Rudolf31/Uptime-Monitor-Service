package models

import "time"

type Note struct {
	Status       string    `json:"status"`
	CheckTime    time.Time `json:"check_time,omitempty"`
	ResponseTime int64     `json:"responce_time"`
}

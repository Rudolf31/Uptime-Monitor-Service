package models

import "time"

type Note struct {
	Status       string    `db:"status"`
	CheckTime    time.Time `db:"check_time"`
	ResponseTime int64     `db:"response_time"`
}

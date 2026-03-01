package models

import "errors"

var (
	ErrInvalidData = errors.New("timeout must be greater than 30 seconds and URL must be valid")
)

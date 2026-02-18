package models

import "errors"

var (
	ErrInvalidData = errors.New("Timeout must be greater than zero")
)

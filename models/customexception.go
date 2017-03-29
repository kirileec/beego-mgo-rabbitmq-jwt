package models

import (
	"time"
)
type CustomException struct {
	Error     string    `json:"error"`
	Exception string    `json:"exception"`
	Message   string    `json:"message"`
	Path      string    `json:"path"`
	TimeStamp time.Time `json:"timestamp"`
}
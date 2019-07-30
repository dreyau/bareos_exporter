package types

import "time"

type LastJob struct {
	Level   string  `json:"level"`
	JobBytes int `json:"job-bytes"`
	JobFiles int `json:"job-files"`
	JobErrors int `json:"job-errors"`
	JobDate time.Time `json:"job-date"`
}

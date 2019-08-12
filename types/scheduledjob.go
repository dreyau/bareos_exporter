package types

import "time"

// ScheduledJob models query result of the time a job is about to be executed
type ScheduledJob struct {
	ScheduledTime time.Time `json:"scheduled-time"`
}

package Types

type LastJob struct {
	Server string `json:"server"`
	Level   string  `json:"level"`
	JobBytes int `json:"job-bytes"`
	JobFiles int `json:"job-files"`
	JobErrors int `json:"job-errors"`
	JobDate string `json:"job-date"`
}

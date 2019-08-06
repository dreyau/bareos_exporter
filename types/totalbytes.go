package types

// TotalBytes models query result of saved bytes sum for a server
type TotalBytes struct {
	Bytes int `json:"files"`
}

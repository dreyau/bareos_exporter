package Queries

import (
	"bareos_exporter/Error"
	"bareos_exporter/Types"
	"database/sql"
	"fmt"
)

type Server struct {
	Name string `json:"name"`
}

func (server Server) TotalBytes(db *sql.DB) Types.TotalBytes {
	query := fmt.Sprintf("SELECT SUM(JobBytes) FROM job WHERE Name='%s'", server.Name)
	results, err := db.Query(query)

	Error.Check(err)

	var totalBytes Types.TotalBytes
	results.Next()
	err = results.Scan(&totalBytes.Bytes)
	results.Close()

	Error.Check(err)

	totalBytes.Server = server.Name

	return totalBytes
}

func (server Server) TotalFiles(db *sql.DB) Types.TotalFiles {
	query := fmt.Sprintf("SELECT SUM(JobFiles) FROM job WHERE Name='%s'", server.Name)
	results, err := db.Query(query)

	Error.Check(err)

	var totalFiles Types.TotalFiles
	results.Next()
	err = results.Scan(&totalFiles.Files)
	results.Close()

	Error.Check(err)

	totalFiles.Server = server.Name

	return totalFiles
}

func (server Server) LastJob(db *sql.DB, lastFullJob bool) Types.LastJob {
	var query string

	if lastFullJob {
		query = fmt.Sprintf("SELECT Level,JobBytes,JobFiles,JobErrors,DATE(StartTime) AS JobDate FROM job WHERE Name LIKE '%s' AND Level = 'F' ORDER BY StartTime DESC LIMIT 1", server.Name)
	} else {
		query = fmt.Sprintf("SELECT Level,JobBytes,JobFiles,JobErrors,DATE(StartTime) AS JobDate FROM job WHERE Name LIKE '%s' ORDER BY StartTime DESC LIMIT 1", server.Name)
	}

	results, err := db.Query(query)

	Error.Check(err)

	var lastJob Types.LastJob
	results.Next()
	err = results.Scan(&lastJob.Level, &lastJob.JobBytes, &lastJob.JobFiles, &lastJob.JobErrors, &lastJob.JobDate)
	results.Close()

	Error.Check(err)

	lastJob.Server = server.Name

	return lastJob
}
package dataaccess

import (
	"bareos_exporter/types"
	"database/sql"
	"fmt"
	"log"
)

type Server struct {
	Name string `json:"name"`
}

func (server Server) TotalBytes(db *sql.DB) (*types.TotalBytes, error) {
	query := fmt.Sprintf("SELECT SUM(JobBytes) FROM job WHERE Name='%s'", server.Name)
	results, err := db.Query(query)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var totalBytes types.TotalBytes
	if results.Next() {
		err = results.Scan(&totalBytes.Bytes)
		results.Close()
	}

	return &totalBytes, err
}

func (server Server) TotalFiles(db *sql.DB) (*types.TotalFiles, error) {
	query := fmt.Sprintf("SELECT SUM(JobFiles) FROM job WHERE Name='%s'", server.Name)
	results, err := db.Query(query)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var totalFiles types.TotalFiles
	if results.Next() {
		err = results.Scan(&totalFiles.Files)
		results.Close()
	}

	return &totalFiles, err
}

func (server Server) LastJob(db *sql.DB, lastFullJob bool) (*types.LastJob, error) {
	var query string

	if lastFullJob {
		query = fmt.Sprintf("SELECT Level,JobBytes,JobFiles,JobErrors,DATE(StartTime) AS JobDate FROM job WHERE Name LIKE '%s' AND Level = 'F' ORDER BY StartTime DESC LIMIT 1", server.Name)
	} else {
		query = fmt.Sprintf("SELECT Level,JobBytes,JobFiles,JobErrors,DATE(StartTime) AS JobDate FROM job WHERE Name LIKE '%s' ORDER BY StartTime DESC LIMIT 1", server.Name)
	}

	results, err := db.Query(query)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var lastJob types.LastJob
	if results.Next() {
		err = results.Scan(&lastJob.Level, &lastJob.JobBytes, &lastJob.JobFiles, &lastJob.JobErrors, &lastJob.JobDate)
		results.Close()
	}

	return &lastJob, err
}
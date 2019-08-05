package dataaccess

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/dreyau/bareos_exporter/types"
	_ "github.com/go-sql-driver/mysql"
)

type Connection struct {
	DB *sql.DB
}

func GetConnection(connectionString string) (*Connection, error) {
	var connection Connection
	var err error

	connection.DB, err = sql.Open("mysql", connectionString)

	return &connection, err
}

func (connection Connection) GetServerList() ([]string, error) {
	date := fmt.Sprintf("%s%", time.Now().Format("2006-01-02"))
	results, err := connection.DB.Query("SELECT DISTINCT Name FROM job WHERE SchedTime LIKE ?", date)

	if err != nil {
		return nil, err
	}

	var servers []string

	for results.Next() {
		var server string
		err = results.Scan(&server)
		servers = append(servers, server)
	}

	return servers, err
}

func (connection Connection) TotalBytes(server string) (*types.TotalBytes, error) {
	results, err := connection.DB.Query("SELECT SUM(JobBytes) FROM job WHERE Name=?", server)

	if err != nil {
		return nil, err
	}

	var totalBytes types.TotalBytes
	if results.Next() {
		err = results.Scan(&totalBytes.Bytes)
		results.Close()
	}

	return &totalBytes, err
}

func (connection Connection) TotalFiles(server string) (*types.TotalFiles, error) {
	results, err := connection.DB.Query("SELECT SUM(JobFiles) FROM job WHERE Name=?", server)

	if err != nil {
		return nil, err
	}

	var totalFiles types.TotalFiles
	if results.Next() {
		err = results.Scan(&totalFiles.Files)
		results.Close()
	}

	return &totalFiles, err
}

func (connection Connection) LastJob(server string) (*types.LastJob, error) {
	results, err := connection.DB.Query("SELECT Level,JobBytes,JobFiles,JobErrors,DATE(StartTime) AS JobDate FROM job WHERE Name LIKE ? ORDER BY StartTime DESC LIMIT 1", server)

	if err != nil {
		return nil, err
	}

	var lastJob types.LastJob
	if results.Next() {
		err = results.Scan(&lastJob.Level, &lastJob.JobBytes, &lastJob.JobFiles, &lastJob.JobErrors, &lastJob.JobDate)
		results.Close()
	}

	return &lastJob, err
}

func (connection Connection) LastFullJob(server string) (*types.LastJob, error) {
	results, err := connection.DB.Query("SELECT Level,JobBytes,JobFiles,JobErrors,DATE(StartTime) AS JobDate FROM job WHERE Name LIKE ? AND Level = 'F' ORDER BY StartTime DESC LIMIT 1", server)

	if err != nil {
		return nil, err
	}

	var lastJob types.LastJob
	if results.Next() {
		err = results.Scan(&lastJob.Level, &lastJob.JobBytes, &lastJob.JobFiles, &lastJob.JobErrors, &lastJob.JobDate)
		results.Close()
	}

	return &lastJob, err
}

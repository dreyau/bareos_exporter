package dataaccess

import (
	"bareos_exporter/types"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Connection struct {
	DB *sql.DB
}

func GetConnection(connectionString string) (Connection, error) {
	var connection Connection
	var err error

	connection.DB, err = sql.Open("mysql", connectionString)

	return connection, err
}

func (connection Connection) GetServerList() ([]string, error) {
	results, err := connection.DB.Query("SELECT DISTINCT Name FROM job WHERE SchedTime LIKE '2019-07-24%'")

	if err != nil{
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

func (connection Connection) TotalBytes(name string) (*types.TotalBytes, error) {
	query := fmt.Sprintf("SELECT SUM(JobBytes) FROM job WHERE Name='%s'", name)
	results, err := connection.DB.Query(query)

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

func (connection Connection) TotalFiles(name string) (*types.TotalFiles, error) {
	query := fmt.Sprintf("SELECT SUM(JobFiles) FROM job WHERE Name='%s'", name)
	results, err := connection.DB.Query(query)

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

func (connection Connection) LastJob(name string) (*types.LastJob, error) {
	query := fmt.Sprintf("SELECT Level,JobBytes,JobFiles,JobErrors,DATE(StartTime) AS JobDate FROM job WHERE Name LIKE '%s' ORDER BY StartTime DESC LIMIT 1", name)

	results, err := connection.DB.Query(query)

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

func (connection Connection) LastFullJob(name string) (*types.LastJob, error) {
	query := fmt.Sprintf("SELECT Level,JobBytes,JobFiles,JobErrors,DATE(StartTime) AS JobDate FROM job WHERE Name LIKE '%s' AND Level = 'F' ORDER BY StartTime DESC LIMIT 1", name)

	results, err := connection.DB.Query(query)

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
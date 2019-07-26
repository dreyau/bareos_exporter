package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"strings"
	"time"
)

type Job struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	ScheduleTime string `json:"schedule-time"`
}

type TotalBytes struct {
	Bytes int   `json:"bytes"`
}

type TotalFiles struct {
	Files int   `json:"files"`
}

type LastJob struct {
	Level   int    `json:"level"`
	JobBytes int `json:"job-bytes"`
	JobFiles int `json:"job-files"`
	JobErrors int `json:"job-errors"`
	JobDate time.Time `json:"job-date"`
}

var (
	mysqlUser     = flag.String("u", "root", "Specify bareos mysql user")
	mysqlAuthFile = flag.String("p", "./auth", "Specify bareos mysql password file path")
	mysqlHostname = flag.String("h", "127.0.0.1", "Specify bareos mysql hostname")
	mysqlPort     = flag.String("P", "3306", "Specify bareos mysql port")
	mysqlDb       = flag.String("db", "bareos", "Specify bareos mysql database name")
)

func init() {
	flag.Usage = func() {
		fmt.Println("Usage: bareos_exporter [ ... ]\n\nParameters:")
		fmt.Println()
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	dat, err := ioutil.ReadFile(*mysqlAuthFile)

	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", *mysqlUser, strings.TrimSpace(string(dat)), *mysqlHostname, *mysqlPort, *mysqlDb)

	db, err := sql.Open("mysql", connection)

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	results, err := db.Query("SELECT DISTINCT JobId, Name, SchedTime FROM job WHERE SchedTime LIKE '2019-07-24%'")

	if err != nil {
		panic(err.Error())
	}

	var jobs []Job

	for results.Next() {
		var job Job
		err = results.Scan(&job.Id, &job.Name, &job.ScheduleTime)

		jobs = append(jobs, job)

		if err != nil {
			panic(err.Error())
		}
	}

	for _, job := range jobs {

		jobBytesQuery := fmt.Sprintf("SELECT SUM(JobBytes) FROM job WHERE Name='%s'", job.Name)
		jobFilesQuery := fmt.Sprintf("SELECT SUM(JobFiles) FROM job WHERE Name='%s'", job.Name)
		lastJobQuery := fmt.Sprintf("SELECT Level,JobBytes,JobFiles,JobErrors,DATE(StartTime) AS JobDate FROM job WHERE Name LIKE '%s' ORDER BY StartTime DESC LIMIT 1", job.Name)
		lastFullJobQuery := fmt.Sprintf("SELECT Level,JobBytes,JobFiles,JobErrors,DATE(StartTime) AS JobDate FROM job WHERE Name LIKE '%s' AND Level = 'F' ORDER BY StartTime DESC LIMIT 1", job.Name)

		jobBytesResults, jobBytesErr := db.Query(jobBytesQuery)
		jobFilesResults, jobFilesErr := db.Query(jobFilesQuery)
		lastJobResults, lastJobErr := db.Query(lastJobQuery)
		lastFullJobResults, lastFullJobErr := db.Query(lastFullJobQuery)

		if jobBytesErr != nil {
			panic(jobBytesErr.Error())
		}
		if jobFilesErr != nil {
			panic(jobFilesErr.Error())
		}
		if lastJobErr != nil {
			panic(lastJobErr.Error())
		}
		if lastFullJobErr != nil {
			panic(lastFullJobErr.Error())
		}

		for jobBytesResults.Next() {
			var totalBytes TotalBytes

			err = results.Scan(&totalBytes.Bytes)

			fmt.Printf("%d bytes saved for server %s\n", totalBytes.Bytes, job.Name)
		}

		for jobFilesResults.Next() {
			var totalFiles TotalFiles

			err = results.Scan(&totalFiles.Files)

			fmt.Printf("%d files saved for server %s\n", totalFiles.Files, job.Name)
		}

		for lastJobResults.Next() {
			var lastJob LastJob

			err = results.Scan(&lastJob.Level, &lastJob.JobBytes, &lastJob.JobFiles, &lastJob.JobErrors, &lastJob.JobDate)

			fmt.Printf("Level %d, %d bytes, %d files, %d errors started on %s for server %s\n", lastJob.Level, lastJob.JobBytes, lastJob.JobFiles, lastJob.JobErrors, lastJob.JobDate, job.Name)
		}

		for lastFullJobResults.Next() {
			var lastJob LastJob

			err = results.Scan(&lastJob.Level, &lastJob.JobBytes, &lastJob.JobFiles, &lastJob.JobErrors, &lastJob.JobDate)

			fmt.Printf("Level %d, %d bytes, %d files, %d errors started on %s for server %s\n", lastJob.Level, lastJob.JobBytes, lastJob.JobFiles, lastJob.JobErrors, lastJob.JobDate, job.Name)
		}
	}
}
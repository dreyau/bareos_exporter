package main

import (
	"bareos_exporter/DataAccess"
	"bareos_exporter/DataAccess/Queries"
	"bareos_exporter/Error"
	"flag"
	"fmt"
)

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

	db := DataAccess.New(*mysqlUser, *mysqlHostname, *mysqlPort, *mysqlDb, *mysqlAuthFile)

	defer db.Close()

	results, err := db.Query("SELECT DISTINCT Name FROM job WHERE SchedTime LIKE '2019-07-22%'")

	Error.Check(err)

	var servers []Queries.Server

	for results.Next() {
		var server Queries.Server
		err = results.Scan(&server.Name)

		servers = append(servers, server)

		Error.Check(err)
	}

	fmt.Println("--- Files per server ---")
	for _, server := range servers {
		files := server.TotalFiles(db)

		fmt.Printf("%s: %d files\n", server.Name, files.Files)
	}

	fmt.Println("--- Bytes per server ---")
	for _, server := range servers {
		bytes := server.TotalBytes(db)

		fmt.Printf("%s: %d bytes\n", server.Name, bytes.Bytes)
	}

	fmt.Println("--- LastJobs ---")
	for _, server := range servers {

		lastJob := server.LastJob(db, false)

		fmt.Printf("%s, %d, %d, %d, %s\n", lastJob.Level, lastJob.JobBytes, lastJob.JobFiles, lastJob.JobErrors, lastJob.JobDate)
	}

	fmt.Println("--- LastFullJobs ---")
	for _, server := range servers {

		lastJob := server.LastJob(db, true)

		fmt.Printf("%s, %d, %d, %d, %s\n", lastJob.Level, lastJob.JobBytes, lastJob.JobFiles, lastJob.JobErrors, lastJob.JobDate)
	}
}
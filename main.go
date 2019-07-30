package main

import (
	"bareos_exporter/dataaccess"
	"bareos_exporter/error"

	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var connection dataaccess.Connection

var (
	mysqlUser     = flag.String("u", "root", "Bareos MySQL username")
	mysqlAuthFile = flag.String("p", "./auth", "MySQL password file path")
	mysqlHostname = flag.String("h", "127.0.0.1", "MySQL hostname")
	mysqlPort     = flag.String("P", "3306", "MySQL port")
	mysqlDb       = flag.String("db", "bareos", "MySQL database name")
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

	pass, err := ioutil.ReadFile(*mysqlAuthFile)
	error.Check(err)

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", *mysqlUser, strings.TrimSpace(string(pass)), *mysqlHostname, *mysqlPort, *mysqlDb)
	db, err := sql.Open("mysql", connectionString)
	error.Check(err)

	connection.DB = db

	defer connection.DB.Close()

	collector := BareosCollector()
	prometheus.MustRegister(collector)

	http.Handle("/metrics", promhttp.Handler())
	log.Info("Beginning to serve on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
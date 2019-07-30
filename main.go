package main

import (
	"bareos_exporter/DataAccess"
	"database/sql"
	"flag"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var db *sql.DB

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

	db = DataAccess.New(*mysqlUser, *mysqlHostname, *mysqlPort, *mysqlDb, *mysqlAuthFile)

	defer db.Close()

	foo := BareosCollector()
	prometheus.MustRegister(foo)

	http.Handle("/metrics", promhttp.Handler())
	log.Info("Beginning to serve on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
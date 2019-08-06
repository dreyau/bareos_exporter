package main

import (
	"github.com/dreyau/bareos_exporter/error"

	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var connectionString string

var (
	exporterPort     = flag.Int("port", 9625, "Bareos exporter port")
	exporterEndpoint = flag.String("endpoint", "/metrics", "Bareos exporter endpoint")
	mysqlUser        = flag.String("u", "root", "Bareos MySQL username")
	mysqlAuthFile    = flag.String("p", "./auth", "MySQL password file path")
	mysqlHostname    = flag.String("h", "127.0.0.1", "MySQL hostname")
	mysqlPort        = flag.String("P", "3306", "MySQL port")
	mysqlDb          = flag.String("db", "bareos", "MySQL database name")
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

	connectionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", *mysqlUser, strings.TrimSpace(string(pass)), *mysqlHostname, *mysqlPort, *mysqlDb)

	collector := bareosCollector()
	prometheus.MustRegister(collector)

	http.Handle(*exporterEndpoint, promhttp.Handler())
	log.Info("Beginning to serve on port ", *exporterPort)

	addr := fmt.Sprintf(":%d", *exporterPort)
	log.Fatal(http.ListenAndServe(addr, nil))
}

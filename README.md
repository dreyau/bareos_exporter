## bareos_exporter
[![Go Report Card](https://goreportcard.com/badge/github.com/dreyau/bareos_exporter)](https://goreportcard.com/report/github.com/dreyau/bareos_exporter)

Prometheus exporter for [bareos](https://github.com/bareos) backup job metrics

### Flags

Name    | Description
--------|-----------------
port    | Bareos exporter port. Defaults to 9625
endpoint| Bareos exporter endpoint. Defaults to "/metrics"
u       | Username used to access Bareos MySQL Database. Defaults to "root"
p       | Path to file containing your MySQL password. Written inside a file to prevent from leaking. Defaults to "./auth"
h       | MySQL instance hostname. Defaults to "127.0.0.1"
P       | MySQL instance port. Defaults to "3306"
db      | MySQL database name. Defaults to "bareos"
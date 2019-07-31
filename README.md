## bareos_exporter

Prometheus exporter for [bareos](https://github.com/bareos) backup job metrics

### Flags

Name    | Description
--------|-----------------
u       | Username used to access Bareos MySQL Database. Defaults to "root"
p       | Path to file containing your MySQL password. Written inside a file to prevent from leaking. Defaults to "./auth"
h       | MySQL instance hostname. Defaults to "127.0.0.1"
P       | MySQL instance port. Defaults to "3306"
db      | MySQL database name. Defaults to "bareos"
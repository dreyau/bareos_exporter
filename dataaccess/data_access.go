package dataaccess

import (
	"bareos_exporter/error"
	"fmt"
	"io/ioutil"
	"strings"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func New(mysqlUser string, mysqlHostname string, mysqlPort string, mysqlDb string, mysqlAuthFile string) *sql.DB{
	pass, err := ioutil.ReadFile(mysqlAuthFile)
	error.Check(err)

	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", mysqlUser, strings.TrimSpace(string(pass)), mysqlHostname, mysqlPort, mysqlDb)
	db, err := sql.Open("mysql", connection)
	error.Check(err)

	return db
}
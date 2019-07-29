package DataAccess

import (
	"bareos_exporter/Error"
	"fmt"
	"io/ioutil"
	"strings"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func New(mysqlUser string, mysqlHostname string, mysqlPort string, mysqlDb string, mysqlAuthFile string) *sql.DB{
	pass, err := ioutil.ReadFile(mysqlAuthFile)

	Error.Check(err)

	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", mysqlUser, strings.TrimSpace(string(pass)), mysqlHostname, mysqlPort, mysqlDb)

	db, err := sql.Open("mysql", connection)

	Error.Check(err)

	return db
}
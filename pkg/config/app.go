package config

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func GetMySQLBase() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "pavel"
	dbPass := "mysqlpaha100688"
	dbName := "testdb2"

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func ErrorCheck(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func PingDB(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	ErrorCheck(err)
}

func DbConnection() (db *sql.DB) {
	db, err := sql.Open("mysql", "testdbwi:testdbwi@tcp(db4free.net:3306)/testdbwi")
	ErrorCheck(err)
	PingDB(db)
	return db
}

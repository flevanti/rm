package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"rm/pkg/refresh"
	"time"
)

var db *sql.DB
var dbUser = "root"
var dbPassword = "root"
var dbHost = "127.0.0.1"
var dbPort = "3306"
var dbName = "rm"

func main() {
	timeStarted := time.Now()
	//ctx, ctxCancelFn := context.WithCancel(context.Background())
	connect()
	defer closeConnection()

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "refresh":
			refresh.Postcodes(db)
		}
	}

	importRM()

	fmt.Printf("Time spent %s\n", time.Since(timeStarted).String())

}

func importRM() {
	fmt.Println("IMPORT COMPLETED")
}

func closeConnection() {
	db.Close()
}

func connect() {
	var err error
	var dbConfig = mysql.Config{
		User:   dbUser,
		Passwd: dbPassword,
		Net:    "tcp",
		Addr:   dbHost + ":" + dbPort,
		DBName: dbName,
	}
	db, err = sql.Open("mysql", dbConfig.FormatDSN())
	if err != nil {
		panic(err)
	}

}

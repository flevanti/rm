package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"rm/pkg/crawl"
	"rm/pkg/refresh"
	"strings"
	"time"
)

var db *sql.DB
var dbUser = "root"
var dbPassword = "root"
var dbHost = "127.0.0.1"
var dbPort = "3306"
var dbName = "rm"

func main() {
	var command string
	timeStarted := time.Now()
	//ctx, ctxCancelFn := context.WithCancel(context.Background())
	connect()
	defer closeConnection()

	if len(os.Args) > 1 {
		command = os.Args[1]
	} else {
		command = "import"
	}
	switch command {
	case "refresh":
		if confirm("Are you sure you want to refresh the postcodes?") {
			refresh.Postcodes(db)
		}
	case "import":
		importRM()
	default:
		fmt.Println("Unknown command")
	}
	fmt.Printf("Time spent %s\n", time.Since(timeStarted).String())
}

func importRM() {
	crawl.Postcode("CB74PL")
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

func confirm(message string) bool {
	rdr := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [y/n]: ", message)

	r, err := rdr.ReadString('\n')
	if err != nil {
		panic(err)
	}

	r = strings.ToLower(strings.TrimSpace(r))

	return r == "y" || r == "yes"

}

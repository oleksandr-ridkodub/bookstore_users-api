package users_db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"

	_ "github.com/go-sql-driver/mysql"
)

var (
	Client   *sql.DB
	username = os.Getenv("MYSQL_USERS_USERNAME")
	password = os.Getenv("MYSQL_USERS_PASSWORD")
	host     = os.Getenv("MYSQL_USERS_HOST")
	schema   = os.Getenv("MYSQL_USERS_SCHEMA")
)

func init() {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		username, password, host, schema,
	)
	var err error
	Client, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}

	if err = Client.Ping(); err != nil {
		panic(err)
	}

	log.Println("Database successfully configured")
}

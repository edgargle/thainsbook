package utils

import (
	"database/sql"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

func ConnectDB() (*sql.DB, error) {
	var db *sql.DB

	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DB_USERNAME")
	cfg.Passwd = os.Getenv("DB_PASSWORD")
	cfg.Net = "tcp"

	if os.Getenv("ENV") == "PROD" {
		cfg.Addr = os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT")
		cfg.TLSConfig = "true"
	} else if os.Getenv("ENV") == "DEV" {
		cfg.Addr = os.Getenv("DB_HOST")
	}

	cfg.DBName = os.Getenv("DB_DATABASE")

	log.Println("Connecting to DB...")
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return nil, pingErr
	}
	log.Println("Successfully connected to MySQL")
	return db, nil
}

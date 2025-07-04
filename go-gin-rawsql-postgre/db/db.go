package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(dsn string) {
	var err error
	for i := 0; i < 10; i++ {
		DB, err = sql.Open("postgres", dsn)
		if err == nil && DB.Ping() == nil {
			log.Println("Connected to DB successfully")
			break
		}
		log.Println("Waiting for database to be ready...")
		time.Sleep(3 * time.Second)
	}
	if DB == nil {
		log.Fatal("Cannot connect to DB after retries")
	}
}

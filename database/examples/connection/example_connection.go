package main

import (
	"log"
	"time"

	"github.com/pureapi/pureapi-core/database"
	"github.com/pureapi/pureapi-core/database/types"

	// Using the SQLite3 driver as an example.
	_ "github.com/mattn/go-sqlite3"
)

// dummyConnectionOpen adapts NewSQLDBAdapter to be used with Connect.
func dummyConnectionOpen(driver, dsn string) (types.DB, error) {
	return database.NewSQLDBAdapter(driver, dsn)
}

func main() {
	cfg := database.ConnectConfig{
		Driver:          "sqlite3",
		Database:        ":memory:", // In-memory SQLite DB.
		ConnMaxLifetime: time.Minute,
		ConnMaxIdleTime: 10 * time.Second,
		MaxOpenConns:    10,
		MaxIdleConns:    5,
	}

	// Connect to the database.
	// For SQLite, the DSN is simply the database name.
	db, err := database.Connect(cfg, dummyConnectionOpen, cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer db.Close()
	log.Println("Database connection established and ping successful.")
}

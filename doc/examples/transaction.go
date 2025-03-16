package main

import (
	"context"
	"fmt"

	"github.com/pureapi/pureapi-core/database"
	"github.com/pureapi/pureapi-core/database/types"
)

func RunDatabase() {
	cfg := database.ConnectConfig{
		Driver:   "sqlite3",
		Database: "example.db",
	}

	db, err := database.Connect(
		cfg, database.NewSQLDBAdapter, "file::memory:?cache=shared",
	)
	if err != nil {
		panic(err)
	}

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	result, err := database.Transaction(
		context.Background(),
		tx,
		func(ctx context.Context, tx types.Tx) (int64, error) {
			_, err := tx.Exec(
				"CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT)",
			)
			if err != nil {
				return 0, err
			}
			res, err := tx.Exec("INSERT INTO users (name) VALUES (?)", "Bob")
			if err != nil {
				return 0, err
			}
			return res.LastInsertId()
		},
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Inserted user with ID:", result)
}

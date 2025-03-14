package examples

import (
	"time"

	"github.com/pureapi/pureapi-core/database"
	"github.com/pureapi/pureapi-core/database/types"
)

// Cfg returns a ConnectConfig for an in-memory SQLite DB.
//
// Returns:
//   - ConnectConfig: The configuration for the database connection.
func Cfg() database.ConnectConfig {
	return database.ConnectConfig{
		Driver:          "sqlite3",
		Database:        ":memory:", // In-memory SQLite DB.
		ConnMaxLifetime: time.Minute,
		ConnMaxIdleTime: 10 * time.Second,
		MaxOpenConns:    10,
		MaxIdleConns:    5,
	}
}

// DummyConnectionOpen adapts NewSQLDBAdapter for connection use.
//
// Parameters:
//   - driver: The database driver name.
//   - dsn: The database connection string.
//
// Returns:
//   - DB: A new instance of DB.
func DummyConnectionOpen(driver string, dsn string) (types.DB, error) {
	return database.NewSQLDBAdapter(driver, dsn)
}

// Connect creates a database connection.
//
// Parameters:
//   - cfg: The configuration for the database connection.
//   - connOpenFn: The function to open the database connection.
//
// Returns:
//   - DB: The database connection.
func Connect(
	cfg database.ConnectConfig, connOpenFn database.ConnOpenFn,
) (types.DB, error) {
	// For SQLite, the DSN is simply the database name.
	db, err := database.Connect(cfg, connOpenFn, cfg.Database)
	if err != nil {
		return nil, err
	}
	return db, nil
}

package database

import (
	"fmt"
	"time"

	"github.com/pureapi/pureapi-core/database/types"
)

// ConnectConfig holds the configuration for the database connection.
type ConnectConfig struct {
	Driver          string        // Driver name. Required to connect.
	User            string        // Database user
	Password        string        // Database password
	Host            string        // Database host
	Port            int           // Database port
	Database        string        // Database name (e.g. "users")
	SocketDirectory string        // Unix socket directory
	SocketName      string        // Unix socket name
	Parameters      string        // Connection parameters
	ConnectionType  string        // Connection type
	ConnMaxLifetime time.Duration // Connection max lifetime
	ConnMaxIdleTime time.Duration // Connection max idle time
	MaxOpenConns    int           // Max open connections
	MaxIdleConns    int           // Max idle connections

	// DSNFormat is an optional format string (e.g. "%s:%s@tcp(%s:%d)/%s?%s").
	// If present (non-empty), it will be used to generate the DSN (with
	// fmt.Sprintf). You can embed placeholders for user, password, host,
	// port, database, and parameters. If left blank, the DSN() function
	// will fall back to a default per-driver build.
	DSNFormat string
}

// ConnOpenFn is a function that opens a database connection.
type ConnOpenFn func(driver string, dsn string) (types.DB, error)

// Connect establishes a connection to the database using the provided
// configuration. It will automatically configure the connection based on the
// provided configuration and then attempt to ping the database.
//
// Parameters:
//   - cfg: The configuration for the database connection.
//   - connOpenFn: The function to open the database connection.
//   - dsn: The database connection string.
//
// Returns:
//   - DB: The database connection.
//   - error: An error if the connection fails.
func Connect(
	cfg ConnectConfig,
	connOpenFn ConnOpenFn,
	dsn string,
) (types.DB, error) {
	db, err := connOpenFn(cfg.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("Connect: failed to open database: %w", err)
	}
	return configureAndPingConnection(db, cfg)
}

// configureAndPingConnection configures the connection and pings the database.
func configureAndPingConnection(
	db types.DB, cfg ConnectConfig,
) (types.DB, error) {
	configureConnection(db, cfg)
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf(
			"configureAndPingConnection: failed to ping database: %w", err,
		)
	}
	return db, nil
}

// configureConnection sets up the runtime connection limits.
func configureConnection(db types.DB, cfg ConnectConfig) {
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
}

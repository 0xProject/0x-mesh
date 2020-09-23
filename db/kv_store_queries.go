// +build !js

package db

import (
	"database/sql"
	"fmt"

	"github.com/ipfs/go-ds-sql"
	_ "github.com/lib/pq" //postgres driver
)

/// PostgreSQL

// Options are the postgres datastore options, reexported here for convenience.
type PostgreSQLOptions struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Table    string
}

type postgreSQLQueries struct {
	tableName string
}

func NewPostgreSQLQueriesForTable(tableName string) *postgreSQLQueries {
	return &postgreSQLQueries{tableName}
}

func (q postgreSQLQueries) Delete() string {
	return `DELETE FROM ` + q.tableName + ` WHERE key = $1`
}

func (q postgreSQLQueries) Exists() string {
	return `SELECT exists(SELECT 1 FROM ` + q.tableName + ` WHERE key=$1)`
}

func (q postgreSQLQueries) Get() string {
	return `SELECT data FROM ` + q.tableName + ` WHERE key = $1`
}

func (q postgreSQLQueries) Put() string {
	return `INSERT INTO ` + q.tableName + ` (key, data) SELECT $1, $2 ON CONFLICT(key) DO UPDATE SET data = $2 WHERE key = $1`
}

func (q postgreSQLQueries) Query() string {
	return `SELECT key, data FROM ` + q.tableName
}

func (q postgreSQLQueries) Prefix() string {
	return ` WHERE key LIKE '%s%%' ORDER BY key`
}

func (q postgreSQLQueries) Limit() string {
	return ` LIMIT %d`
}

func (q postgreSQLQueries) Offset() string {
	return ` OFFSET %d`
}

func (q postgreSQLQueries) GetSize() string {
	return `SELECT octet_length(data) FROM ` + q.tableName + ` WHERE key = $1`
}

// Create returns a datastore connected to postgres initialized with a table
func (opts *PostgreSQLOptions) CreatePostgres() (*sqlds.Datastore, error) {
	opts.setDefaults()
	fmtstr := "postgresql:///%s?host=%s&port=%s&user=%s&password=%s&sslmode=disable"
	constr := fmt.Sprintf(fmtstr, opts.Database, opts.Host, opts.Port, opts.User, opts.Password)
	db, err := sql.Open("postgres", constr)
	if err != nil {
		return nil, err
	}

	createTable := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (key TEXT NOT NULL UNIQUE, data BYTEA NOT NULL)", opts.Table)
	_, err = db.Exec(createTable)

	if err != nil {
		return nil, err
	}

	return sqlds.NewDatastore(db, NewPostgreSQLQueriesForTable(opts.Table)), nil
}

func (opts *PostgreSQLOptions) setDefaults() {
	if opts.Table == "" {
		opts.Table = "kv"
	}
	if opts.Host == "" {
		opts.Host = "postgres"
	}

	if opts.Port == "" {
		opts.Port = "5432"
	}

	if opts.User == "" {
		opts.User = "postgres"
	}

	if opts.Database == "" {
		opts.Database = "datastore"
	}
}

/// Sqlite

type sqliteQueries struct {
	tableName string
}

func NewSqliteQueriesForTable(tableName string) *sqliteQueries {
	return &sqliteQueries{tableName}
}

func (q sqliteQueries) Delete() string {
	return `DELETE FROM ` + q.tableName + ` WHERE key = $1`
}

func (q sqliteQueries) Exists() string {
	return `SELECT exists(SELECT 1 FROM ` + q.tableName + ` WHERE key=$1)`
}

func (q sqliteQueries) Get() string {
	return `SELECT data FROM ` + q.tableName + ` WHERE key = $1`
}

func (q sqliteQueries) Put() string {
	return `INSERT INTO ` + q.tableName + ` (key, data) SELECT $1, $2 ON CONFLICT(key) DO UPDATE SET data = $2 WHERE key = $1`
}

func (q sqliteQueries) Query() string {
	return `SELECT key, data FROM ` + q.tableName
}

func (q sqliteQueries) Prefix() string {
	return ` WHERE key LIKE '%s%%' ORDER BY key`
}

func (q sqliteQueries) Limit() string {
	return ` LIMIT %d`
}

func (q sqliteQueries) Offset() string {
	return ` OFFSET %d`
}

func (q sqliteQueries) GetSize() string {
	return `SELECT length(data) FROM ` + q.tableName + ` WHERE key = $1`
}

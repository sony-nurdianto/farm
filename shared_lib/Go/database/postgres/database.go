package postgres

import (
	"context"
	"database/sql"
	"time"
)

// type Stmt interface {
// 	QueryRowContext(ctx context.Context, args ...any) Row
// 	QueryContext(ctx context.Context, args ...any) (Rows, error)
// }
//
// type Rows interface {
// 	Close() error
// 	Next() bool
// 	Scan(dest ...any) error
// }
//
// type Row interface {
// 	Scan(dest ...any) error
// }

type PostgresDatabase interface {
	SetMaxOpenConns(n int)
	SetMaxIdleConns(n int)
	SetConnMaxLifetime(d time.Duration)
	PingContext(ctx context.Context) error

	Prepare(query string) (*sql.Stmt, error)
	// PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	// ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)

	Close() error
}

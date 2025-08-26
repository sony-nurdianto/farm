package pkg

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

//go:generate mockgen -package=mocks -destination=../test/mocks/mock_postgresdb.go -source=database_def.go
type PostgresDatabase interface {
	SetMaxOpenConns(n int)
	SetMaxIdleConns(n int)
	SetConnMaxLifetime(d time.Duration)
	PingContext(ctx context.Context) error

	Begin() (SQLTx, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (SQLTx, error)

	Prepare(query string) (Stmt, error)
	Close() error
}

type postgresDatabase struct {
	db *sql.DB
}

func NewPostgresDb(db *sql.DB) PostgresDatabase {
	return postgresDatabase{db}
}

func (pdb postgresDatabase) SetMaxOpenConns(n int) {
	pdb.db.SetMaxOpenConns(n)
}

func (pdb postgresDatabase) SetMaxIdleConns(n int) {
	pdb.db.SetMaxIdleConns(n)
}

func (pdb postgresDatabase) SetConnMaxLifetime(d time.Duration) {
	pdb.db.SetConnMaxLifetime(d)
}

func (pdb postgresDatabase) PingContext(ctx context.Context) error {
	return pdb.db.PingContext(ctx)
}

func (pdb postgresDatabase) Prepare(query string) (Stmt, error) {
	stmt, err := pdb.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	return NewStmt(stmt), nil
}

func (pdb postgresDatabase) Begin() (SQLTx, error) {
	tx, err := pdb.db.Begin()
	if err != nil {
		return nil, err
	}
	return NewSQLTx(tx), nil
}

func (pdb postgresDatabase) BeginTx(ctx context.Context, opts *sql.TxOptions) (SQLTx, error) {
	tx, err := pdb.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}

	return NewSQLTx(tx), nil
}

func (pdb postgresDatabase) Close() error {
	return pdb.db.Close()
}

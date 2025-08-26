package pkg

import (
	"context"
	"database/sql"
)

//go:generate mockgen -package=mocks -destination=../test/mocks/mock_pgstmt.go -source=stmt_def.go
type Stmt interface {
	QueryRowContext(ctx context.Context, args ...any) Row
	QueryContext(ctx context.Context, args ...any) (Rows, error)
	Close() error
	ToSQLSTMT() *sql.Stmt
}

type stmt struct {
	statement *sql.Stmt
}

func NewStmt(s *sql.Stmt) stmt {
	return stmt{statement: s}
}

func (s stmt) QueryRowContext(ctx context.Context, args ...any) Row {
	return s.statement.QueryRowContext(ctx, args...)
}

func (s stmt) QueryContext(ctx context.Context, args ...any) (Rows, error) {
	return s.statement.QueryContext(ctx, args...)
}

func (s stmt) ToSQLSTMT() *sql.Stmt {
	return s.statement
}

func (s stmt) Close() error {
	return s.statement.Close()
}

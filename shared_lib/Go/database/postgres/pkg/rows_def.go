package pkg

//go:generate mockgen -package=mocks -destination=../test/mocks/mock_pgrows.go -source=rows_def.go
type Rows interface {
	Close() error
	Next() bool
	Scan(dest ...any) error
}

package pkg

//go:generate mockgen -package=mocks -destination=../test/mocks/mock_pgrow.go -source=row_def.go
type Row interface {
	Scan(dest ...any) error
	Err() error
}

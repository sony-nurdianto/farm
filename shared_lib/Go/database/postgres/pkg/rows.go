package pkg

type Rows interface {
	Close() error
	Next() bool
	Scan(dest ...any) error
}

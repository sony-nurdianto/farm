package pkg

type Row interface {
	Scan(dest ...any) error
}

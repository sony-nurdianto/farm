package pkg

import "database/sql"

type PostgresInstance interface {
	Open(driverName string, dataSourceName string) (PostgresDatabase, error)
}

type pgi struct{}

func NewPostgresInstance() pgi {
	return pgi{}
}

func (i pgi) Open(driverName, dataSourceName string) (PostgresDatabase, error) {
	sql, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	return NewPostgresDb(sql), nil
}

package postgres

import (
	"context"
	"time"

	_ "github.com/lib/pq"
)

type postgres struct {
	db PostgresDatabase
}

func OpenPostgres(uri string, instance PostgresInstance) (
	PostgresDatabase, error,
) {
	db, err := instance.Open("postgres", uri)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(30)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

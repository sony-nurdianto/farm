package pkg

import (
	"context"
	"time"

	_ "github.com/lib/pq"
)

func OpenPostgres(uri string, instance PostgresInstance) (
	pg PostgresDatabase, _ error,
) {
	db, err := instance.Open("postgres", uri)
	if err != nil {
		return pg, err
	}

	db.SetMaxOpenConns(30)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return pg, err
	}

	return db, nil
}

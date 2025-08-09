package integrations_test

import (
	"testing"

	"github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres/pkg"
)

func TestOpenPostgres(t *testing.T) {
	instance := pkg.NewPostgresInstance()
	db, err := pkg.OpenPostgres("postgres://sony:secret@localhost:5000/test?sslmode=disable", instance)
	if err != nil {
		t.Fatalf("Expected Success Open Connection To Postgres Database but got %s", err)
	}

	if db == nil {
		t.Fatalf("Expected OpenPostgres Method return non nil value")
	}

	_, ok := db.(pkg.PostgresDatabase)
	if !ok {
	}
}

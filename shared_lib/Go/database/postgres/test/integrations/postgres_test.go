package integrations_test

import (
	"testing"

	"github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres"
)

func TestOpenPostgres(t *testing.T) {
	pgi := postgres.NewPostgresInstance()

	db, err := postgres.OpenPostgres("postgres://sony:secret@localhost:5000/farmer_db?sslmode=disable", pgi)
	if err != nil {
		t.Errorf("Failed to opendb: %s", err)
	}

	if db == nil {
		t.Fatalf("Expected OpenPostgres Method return non nil value")
	}

	db.Close()
}

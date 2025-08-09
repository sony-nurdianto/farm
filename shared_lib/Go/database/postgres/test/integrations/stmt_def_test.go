package integrations_test

import (
	"context"
	"testing"

	"github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres/pkg"
	"github.com/stretchr/testify/assert"
)

func TestStmtDef(t *testing.T) {
	t.Run("QueryContext Error", func(t *testing.T) {
		ins := pkg.NewPostgresInstance()

		db, err := ins.Open("postgres", "postgres://sony:secret@localhost:5000/test?sslmode=disable")
		assert.NoError(t, err)

		query := `
			insert into users (name,message) values ($1,$2) returning name, message
		`

		stmt, err := db.Prepare(query)
		assert.NoError(t, err)

		_, err = stmt.QueryContext(context.Background(), "something funny", "Slim Shady", "Whats up")
		assert.NoError(t, err)
	})

	t.Run("QueryRowContext", func(t *testing.T) {
		ins := pkg.NewPostgresInstance()

		db, err := ins.Open("postgres", "postgres://sony:secret@localhost:5000/test?sslmode=disable")
		assert.NoError(t, err)

		query := `
			insert into users (name,message) values ($1,$2) returning name, message
		`

		stmt, err := db.Prepare(query)
		assert.NoError(t, err)

		_ = stmt.QueryRowContext(context.Background(), "Slim Shady", "Hallo My Name Is ?")
	})
}

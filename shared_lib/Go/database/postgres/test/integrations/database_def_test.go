package integrations_test

import (
	"testing"

	"github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres/pkg"
	"github.com/stretchr/testify/assert"
)

func TestPostgresDatabase(t *testing.T) {
	t.Run("Prepare Test", func(t *testing.T) {
		ins := pkg.NewPostgresInstance()

		db, err := ins.Open("postgres", "postgres://sony:secret@localhost:5000/test?sslmode=disable")
		assert.NoError(t, err)

		defer db.Close()

		query := `
			insert into users (name,message) values ($1,$2) returning name, message
		`

		_, err = db.Prepare(query)
		assert.NoError(t, err)
	})

	t.Run("Prepare Test Error", func(t *testing.T) {
		ins := pkg.NewPostgresInstance()

		db, err := ins.Open("postgres", "postgres://sony:secret@localhost:5000/test?sslmode=disable")
		assert.NoError(t, err)

		defer db.Close()

		query := `
			hmms something error
		`

		_, err = db.Prepare(query)
		assert.Error(t, err)
	})
}

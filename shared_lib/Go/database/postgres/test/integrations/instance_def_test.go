package integrations_test

import (
	"testing"

	"github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres/pkg"
	"github.com/stretchr/testify/assert"
)

func TestPostgresInstance(t *testing.T) {
	t.Run("Open", func(t *testing.T) {
		ins := pkg.NewPostgresInstance()

		_, err := ins.Open("postgres", "postgres://sony:secret@localhost:5000/test?sslmode=disable")
		assert.NoError(t, err)
	})

	t.Run("Open Error wrong driver", func(t *testing.T) {
		ins := pkg.NewPostgresInstance()

		_, err := ins.Open("something", "postgres://sony:secret@localhost:5000/test?sslmode=disable")
		assert.Error(t, err)
	})

	t.Run("Open Error wrong address", func(t *testing.T) {
		ins := pkg.NewPostgresInstance()

		_, err := ins.Open("postgres", "something")
		assert.Error(t, err)
	})
}

package repo_test

import (
	"testing"

	"github.com/sony-nurdianto/farm/auth/internal/repository"
	"github.com/sony-nurdianto/farm/shared_lib/Go/mykafka/pkg"
)

func TestRepo_CreateNewRepo(t *testing.T) {
	registery := pkg.NewRegistery()
	repo, err := repository.NewPostgresRepo(registery)
	if err != nil {
		t.Errorf("Expected NewPostgresRepo Method did not return error but got %s", err)
	}

	t.Run("Create New User Test", func(t *testing.T) {
		email := "test@email.com"
		password := "this is password hash"
		user, err := repo.CreateUser(email, password)
		if err != nil {
			t.Errorf("Failed create users %s", err)
		}

		if user.Email != email {
			t.Errorf("Expected email is %s but got %s", email, user.Email)
		}

		if user.Password != password {
			t.Errorf("Expected password is %s but got %s", password, user.Password)
		}
	})
}

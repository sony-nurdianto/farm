package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/sony-nurdianto/farm/auth/internal/constants"
	"github.com/sony-nurdianto/farm/auth/internal/entity"
	"github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres/pkg"
)

type repoPostgres struct {
	// db             pkg.PostgresDatabase
	createUserstmt pkg.Stmt
}

func prepareStmt(query string, db pkg.PostgresDatabase) (pkg.Stmt, error) {
	facQuery := fmt.Sprintf(
		query,
		constants.USERS_TABLE,
	)

	return db.Prepare(facQuery)
}

func NewPostgresRepo() (rp repoPostgres, _ error) {
	pgi := pkg.NewPostgresInstance()

	db, err := pkg.OpenPostgres("postgres://sony:secret@localhost:5000/auth?sslmode=disable", pgi)
	if err != nil {
		return rp, err
	}

	crs, err := prepareStmt(constants.QUERY_CREATE_USERS, db)
	if err != nil {
		return rp, err
	}

	rp.createUserstmt = crs

	return rp, nil
}

func (rp repoPostgres) CreateUser(email, password string) (user entity.Users, _ error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Millisecond*500,
	)
	defer cancel()

	row := rp.createUserstmt.QueryRowContext(ctx, email, password)

	err := row.Scan(&user.Id, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, err
	}

	return user, nil
}

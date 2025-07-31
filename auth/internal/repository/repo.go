package repository

import (
	"database/sql"
	"fmt"

	"github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres"
)

type repoPostgres struct {
	db                 postgres.PostgresDatabase
	registerFarmerStmt *sql.Stmt
}

func prepareStmt(queryStr string, db postgres.PostgresDatabase) (*sql.Stmt, error) {
	facQuery := fmt.Sprintf(
		queryStr,
		"table name",
	)

	return db.Prepare(facQuery)
}

func NewRepoPostgres(uri string, pgi postgres.PostgresInstance) (rp repoPostgres, _ error) {
	instance := postgres.NewPostgresInstance()

	db, err := postgres.OpenPostgres(uri, instance)
	if err != nil {
		return rp, err
	}

	rp.db = db

	farmerStmt, err := prepareStmt("", db)
	if err != nil {
		return rp, err
	}

	rp.registerFarmerStmt = farmerStmt

	return rp, nil
}

func (rp repoPostgres) RegisterFarmer() {
	stmt, err := rp.db.Prepare("")
	if err != nil {
		panic(err)
	}

	stmt.QueryRow("")
}

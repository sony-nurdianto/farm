package repository

// import (
// 	"context"
// 	"database/sql"
// 	"fmt"
// 	"time"
//
// 	"github.com/sony-nurdianto/farm/auth/internal/constants"
// 	"github.com/sony-nurdianto/farm/auth/internal/entity"
// 	"github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres"
// )
//
// type repoPostgres struct {
// 	db                 postgres.PostgresDatabase
// 	registerFarmerStmt *sql.Stmt
// }
//
// func prepareStmt(queryStr string, db postgres.PostgresDatabase) (*sql.Stmt, error) {
// 	facQuery := fmt.Sprintf(
// 		queryStr,
// 		constants.USERS_TABLE,
// 	)
//
// 	return db.Prepare(facQuery)
// }
//
// func NewRepoPostgres(uri string, pgi postgres.PostgresInstance) (rp repoPostgres, _ error) {
// 	instance := postgres.NewPostgresInstance()
//
// 	db, err := postgres.OpenPostgres(uri, instance)
// 	if err != nil {
// 		return rp, err
// 	}
//
// 	rp.db = db
//
// 	farmerStmt, err := prepareStmt(constants.QUERY_CREATE_USERS, db)
// 	if err != nil {
// 		return rp, err
// 	}
//
// 	rp.registerFarmerStmt = farmerStmt
//
// 	return rp, nil
// }
//
// func (rp repoPostgres) RegisterFarmer(email, password string) (user entity.Users, _ error) {
// 	ctx, cancel := context.WithTimeout(
// 		context.Background(),
// 		time.Millisecond*500,
// 	)
//
// 	defer cancel()
// 	row := rp.registerFarmerStmt.QueryRowContext(ctx, "email", "password")
//
// 	err := row.Scan(&user.Id, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
// 	if err != nil {
// 		return user, err
// 	}
//
// 	return user, nil
// }

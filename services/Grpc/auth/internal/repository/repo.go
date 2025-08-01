package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/sony-nurdianto/farm/auth/internal/constants"
	"github.com/sony-nurdianto/farm/auth/internal/entity"
	"github.com/sony-nurdianto/farm/auth/internal/models"
	"github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres/pkg"
	schema "github.com/sony-nurdianto/farm/shared_lib/Go/mykafka/pkg"
)

type RepoPostgres struct {
	schemaRegistery    schema.SchemaRegistery
	createUserstmt     pkg.Stmt
	getUserByEmailStmt pkg.Stmt
}

func prepareStmt(query string, db pkg.PostgresDatabase) (pkg.Stmt, error) {
	facQuery := fmt.Sprintf(
		query,
		constants.USERS_TABLE,
	)

	return db.Prepare(facQuery)
}

func NewPostgresRepo() (rp RepoPostgres, _ error) {
	srgs, err := schema.NewSchemaRegistery("http://localhost:8081")
	if err != nil {
		return rp, err
	}

	rp.schemaRegistery = srgs

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

	gue, err := prepareStmt(constants.QUERY_GET_USER_BY_EMAIL, db)
	if err != nil {
		return rp, err
	}

	rp.getUserByEmailStmt = gue

	return rp, nil
}

func createUserSchema(name, sch string) {
}

func (rp RepoPostgres) CreateUserAsync(email, passwordHash string, schemaVersion int) error {
	user := &models.InsertUser{
		Id:       uuid.NewString(),
		Email:    email,
		Password: passwordHash,
	}

	schemaName := "users"

	_, err := rp.schemaRegistery.GetSchemaRegistery(schemaName, schemaVersion)
	if err != nil {
		if errors.Is(err, schema.SchemaIsNotFoundErr) {

			_, err = rp.schemaRegistery.RegisterSchema(schemaName, user.Schema(), false)
			if err != nil {
				return err
			}
		}

		return err
	}

	av, err := schema.NewAvroGenericSerde(rp.schemaRegistery.Client())
	if err != nil {
		return err
	}

	_, err = av.Serialize(schemaName, user)
	if err != nil {
		log.Fatalln(err)
	}

	return nil

	// producer := schema.NewKafkaProducer()
	// producer.Producer()
}

func (rp RepoPostgres) CreateUser(email, passwordHash string) (user entity.Users, _ error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Millisecond*500,
	)
	defer cancel()

	userId := uuid.NewString()

	row := rp.createUserstmt.QueryRowContext(ctx, userId, email, passwordHash)

	err := row.Scan(&user.Id, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (rp RepoPostgres) GetUserByEmail(email string) (user entity.Users, _ error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Millisecond*500,
	)

	defer cancel()

	row := rp.getUserByEmailStmt.QueryRowContext(ctx, email)

	err := row.Scan(&user.Id, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, err
	}

	return user, nil
}

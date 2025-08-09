package repository

import (
	"fmt"

	"github.com/sony-nurdianto/farm/auth/internal/constants"
	"github.com/sony-nurdianto/farm/auth/internal/entity"
	"github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres/pkg"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/avr"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/kev"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/schrgs"
)

//go:generate mockgen -package=mocks -destination=../../test/mocks/mock_schrgs.go  github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/schrgs SchemaRegisteryClient,SchemaRegisteryInstance
//go:generate mockgen -package=mocks -destination=../../test/mocks/mock_postgres.go  github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres/pkg PostgresInstance,PostgresDatabase,Stmt,Row
//go:generate mockgen -destination=../../test/mocks/mock_confluent_client.go -package=mocks github.com/confluentinc/confluent-kafka-go/v2/schemaregistry Client
//go:generate mockgen -package=mocks -destination=../../test/mocks/mock_avr.go  github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/avr AvrSerdeInstance,AvrSerializer,AvrDeserializer
//go:generate mockgen -package=mocks -destination=../../test/mocks/mock_kev.go  github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/kev Kafka,KevProducer
//go:generate mockgen -package=mocks -destination=../../test/mocks/mock_authrepo.go -source=repo.go

type AuthRepo interface {
	CreateUserAsync(id, email, fullName, phone, passwordHash string) error
	GetUserByEmail(email string) (user entity.Users, _ error)
}

type authRepo struct {
	schemaRegistery       *schrgs.SchemaRegistery
	schemaRegisteryClient schrgs.SchemaRegisteryClient
	avro                  avr.AvrSerdeInstance
	kafka                 *kev.KafkaProducerPool
	createUserstmt        pkg.Stmt
	getUserByEmailStmt    pkg.Stmt
}

func prepareStmt(query string, db pkg.PostgresDatabase) (pkg.Stmt, error) {
	facQuery := fmt.Sprintf(
		query,
		constants.ACCOUNT_TABLE,
	)

	return db.Prepare(facQuery)
}

func NewAuthRepo(
	sri schrgs.SchemaRegisteryInstance,
	pgi pkg.PostgresInstance,
	avr avr.AvrSerdeInstance,
	kv kev.Kafka,
) (rp authRepo, _ error) {
	srgs, err := schrgs.NewSchemaRegistery("http://localhost:8081", sri)
	if err != nil {
		return rp, err
	}

	rp.schemaRegistery = &srgs
	rp.schemaRegisteryClient = srgs.Client()

	rp.avro = avr

	pool := kev.NewKafkaProducerPool(kv, nil)
	rp.kafka = pool

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

// func (rp AuthRepo) CreateUser(email, passwordHash string) (user entity.Users, _ error) {
// 	ctx, cancel := context.WithTimeout(
// 		context.Background(),
// 		time.Millisecond*500,
// 	)
// 	defer cancel()
//
// 	userId := uuid.NewString()
//
// 	row := rp.createUserstmt.QueryRowContext(ctx, userId, email, passwordHash)
//
// 	err := row.Scan(&user.Id, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
// 	if err != nil {
// 		return user, err
// 	}
//
// 	return user, nil
// }
//

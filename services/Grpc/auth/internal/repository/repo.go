package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sony-nurdianto/farm/auth/internal/constants"
	"github.com/sony-nurdianto/farm/auth/internal/entity"
	"github.com/sony-nurdianto/farm/auth/internal/models"
	"github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres/pkg"
	kk "github.com/sony-nurdianto/farm/shared_lib/Go/mykafka/pkg"
)

type RepoPostgres struct {
	schemaRegistery    kk.SchemaRegistery
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
	srgs, err := kk.NewSchemaRegistery("http://localhost:8081")
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

func ensureSchema(
	subject string,
	version int,
	schema string,
	topic string,
	msg any,
	registry kk.SchemaRegistery,
) error {
	_, err := registry.GetSchemaRegistery(subject, version)
	if err == nil {
		return nil
	}

	if !errors.Is(err, kk.SchemaIsNotFoundErr) {
		return err
	}

	_, regErr := registry.RegisterSchema(subject, schema, false)
	if regErr != nil {
		return regErr
	}

	return publishAvro(registry, topic, msg)
}

func publishAvro(client kk.SchemaRegistery, topic string, msg any) error {
	av, err := kk.NewAvroGenericSerde(client.Client())
	if err != nil {
		return err
	}

	payload, err := av.Serialize(topic, msg)
	if err != nil {
		return err
	}

	cfg := map[kk.ConfigKeyKafka]string{
		kk.BOOTSTRAP_SERVERS:                     "localhost:29092",
		kk.ACKS:                                  "all",
		kk.ENABLE_IDEMPOTENCE:                    "true",
		kk.COMPRESION_TYPE:                       "snappy",
		kk.RETRIES:                               "5",
		kk.RETRY_BACKOFF_MS:                      "100",
		kk.LINGER_MS:                             "5",
		kk.MAX_IN_FLIGHT_REQUESTS_PER_CONNECTION: "5",
	}

	record := &kk.MessageKafka{
		TopicPartition: kk.KafkaTopicPartition{
			Topic:     &topic,
			Partition: kk.KafkaPartitionAny,
		},
		Value: payload,
	}

	return kk.NewKafkaProducerPool().Produce(cfg, record)
}

func (rp RepoPostgres) CreateUserAsync(id, email, passwordHash string, schemaVersion int) error {
	user := &models.InsertUser{
		Id:       id,
		Email:    email,
		Password: passwordHash,
	}

	topic := "insert-users"
	subject := topic + "-value"

	if err := ensureSchema(subject, schemaVersion, user.Schema(), topic, user, rp.schemaRegistery); err != nil {
		return err
	}

	return publishAvro(rp.schemaRegistery, topic, user)
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

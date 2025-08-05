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
	schemaRegistery    kk.RegisterySchema
	createUserstmt     pkg.Stmt
	getUserByEmailStmt pkg.Stmt
}

func prepareStmt(query string, db pkg.PostgresDatabase) (pkg.Stmt, error) {
	facQuery := fmt.Sprintf(
		query,
		constants.ACCOUNT_TABLE,
	)

	return db.Prepare(facQuery)
}

func NewPostgresRepo(rgs kk.SchemaRegistery) (rp RepoPostgres, _ error) {
	srgs, err := kk.NewSchemaRegistery("http://localhost:8081", rgs)
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
	schema string,
	registry kk.RegisterySchema,
) error {
	_, err := registry.GetLatestSchemaRegistery(subject)
	if err == nil {
		return nil
	}

	if !errors.Is(err, kk.SchemaIsNotFoundErr) {
		return err
	}

	_, regErr := registry.CreateAvroSchema(subject, schema, false)
	if regErr != nil {
		return regErr
	}

	return nil
}

func publishAvro(
	client kk.RegisterySchema,
	accountTopic string,
	userTopic string,
	account *models.InsertAccount,
	user *models.InsertFarmerUser,
) error {
	av, err := kk.NewAvroGenericSerde(client.Client())
	if err != nil {
		return err
	}

	accountPayload, err := av.Serialize(accountTopic, account)
	if err != nil {
		return err
	}

	userPayload, err := av.Serialize(userTopic, user)
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

		kk.TRANSACTIONAL_ID: "register-user",
	}

	accountRecord := kk.MessageKafka{
		TopicPartition: kk.KafkaTopicPartition{
			Topic:     &accountTopic,
			Partition: kk.KafkaPartitionAny,
		},
		Value: accountPayload,
	}.Factory()

	userRecord := kk.MessageKafka{
		TopicPartition: kk.KafkaTopicPartition{
			Topic:     &userTopic,
			Partition: kk.KafkaPartitionAny,
		},
		Value: userPayload,
	}.Factory()

	pool := kk.NewKafkaProducerPool()
	producer, err := pool.Producer(cfg)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Millisecond*1000,
	)
	defer cancel()

	if err := producer.InitTransactions(ctx); err != nil {
		return err
	}

	if err := producer.BeginTransaction(); err != nil {
		return err
	}

	if err := producer.Produce(&accountRecord, nil); err != nil {
		return errors.Join(err, producer.AbortTransaction(ctx))
	}

	if err := producer.Produce(&userRecord, nil); err != nil {
		return errors.Join(err, producer.AbortTransaction(ctx))
	}

	return producer.CommitTransaction(ctx)
}

func (rp RepoPostgres) CreateUserAsync(
	id, email, fullName, phone, passwordHash string,
) error {
	account := &models.InsertAccount{
		Id:       id,
		Email:    email,
		Password: passwordHash,
	}

	user := &models.InsertFarmerUser{
		Id:       id,
		FullName: fullName,
		Email:    email,
		Phone:    phone,
	}

	valueStr := "-value"
	accountTopic := "insert-account"
	accountSubject := fmt.Sprintf("%s-%s", accountTopic, valueStr)

	userTopic := "insert-user"
	userSubject := fmt.Sprintf("%s-%s", userTopic, valueStr)

	err := ensureSchema(accountSubject, account.Schema(), rp.schemaRegistery)
	if err != nil {
		return err
	}

	err = ensureSchema(userSubject, user.Schema(), rp.schemaRegistery)
	if err != nil {
		return err
	}

	return publishAvro(rp.schemaRegistery, accountTopic, userTopic, account, user)
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

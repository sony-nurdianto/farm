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
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/avr"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/kev"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/schrgs"
)

type RepoPostgres struct {
	schemaRegistery       *schrgs.SchemaRegistery
	schemaRegisteryClient schrgs.SchemaRegisteryClient
	avro                  avr.AvrSerdeInstance
	kafka                 kev.Kafka
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

func NewPostgresRepo(sri schrgs.SchemaRegisteryInstance) (rp RepoPostgres, _ error) {
	srgs, err := schrgs.NewSchemaRegistery("http://localhost:8081", sri)
	if err != nil {
		return rp, err
	}
	rp.schemaRegistery = &srgs
	rp.schemaRegisteryClient = srgs.Client()

	avr := avr.NewAvrSerdeInstance()
	rp.avro = avr

	kk := kev.NewKafka()
	rp.kafka = kk

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

func (rp RepoPostgres) ensureSchema(
	subject string,
	schema string,
) error {
	_, err := rp.schemaRegistery.GetLatestSchemaRegistery(subject)
	if err == nil {
		return nil
	}

	if !errors.Is(err, schrgs.SchemaIsNotFoundErr) {
		return err
	}

	_, regErr := rp.schemaRegistery.CreateAvroSchema(subject, schema, false)
	if regErr != nil {
		return regErr
	}

	return nil
}

func (rp RepoPostgres) publishAvro(
	accountTopic string,
	userTopic string,
	account *models.InsertAccount,
	user *models.InsertFarmerUser,
) error {
	av, err := avr.NewAvroGenericSerde(rp.schemaRegisteryClient.Client(), rp.avro)
	if err != nil {
		return err
	}

	accountPayload, err := av.Serializer.Serialize(accountTopic, account)
	if err != nil {
		return err
	}

	userPayload, err := av.Serializer.Serialize(userTopic, user)
	if err != nil {
		return err
	}

	cfg := map[kev.ConfigKeyKafka]string{
		kev.BOOTSTRAP_SERVERS:  "localhost:29092",
		kev.ACKS:               "all",
		kev.ENABLE_IDEMPOTENCE: "true",
		kev.COMPRESION_TYPE:    "snappy",
		kev.RETRIES:            "5",
		kev.RETRY_BACKOFF_MS:   "100",
		kev.LINGER_MS:          "5",
		kev.MAX_IN_FLIGHT_REQUESTS_PER_CONNECTION: "5",

		kev.TRANSACTIONAL_ID: "register-user",
	}

	accountRecord := kev.MessageKafka{
		TopicPartition: kev.KafkaTopicPartition{
			Topic:     &accountTopic,
			Partition: kev.KafkaPartitionAny,
		},
		Value: accountPayload,
	}.Factory()

	userRecord := kev.MessageKafka{
		TopicPartition: kev.KafkaTopicPartition{
			Topic:     &userTopic,
			Partition: kev.KafkaPartitionAny,
		},
		Value: userPayload,
	}.Factory()

	pool := kev.NewKafkaProducerPool(rp.kafka, nil)
	producer, err := pool.Producer(cfg)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Millisecond*1000,
	)
	defer cancel()

	for i := range 5 {
		err := producer.InitTransactions(ctx)
		if err == nil {
			fmt.Printf("✅ Kafka transactional producer ready after %d attempt(s)\n", i)
			break
		}
		log.Println("Waiting Init Transaction Ready")
		time.Sleep(2 * time.Second)
	}

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

	err := rp.ensureSchema(accountSubject, account.Schema())
	if err != nil {
		return err
	}

	err = rp.ensureSchema(userSubject, user.Schema())
	if err != nil {
		return err
	}

	return rp.publishAvro(accountTopic, userTopic, account, user)
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

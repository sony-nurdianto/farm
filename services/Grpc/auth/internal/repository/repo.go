package repository

import (
	"context"
	"fmt"
	"os"
	"time"

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

const TRANSACTIONAL_ID string = "register-user"

type AuthRepo interface {
	CreateUserAsync(ctx context.Context, id, email, fullName, phone, passwordHash string) error
	GetUserByEmail(ctx context.Context, email string) (user entity.Users, _ error)
}

type authRepo struct {
	schemaRegisteryClient schrgs.SchemaRegisteryClient
	avroSerializer        avr.AvrSerializer
	db                    pkg.PostgresDatabase
	createUserstmt        pkg.Stmt
	getUserByEmailStmt    pkg.Stmt
	authProducer          kev.KevProducer
}

func prepareStmt(query string, db pkg.PostgresDatabase) (pkg.Stmt, error) {
	facQuery := fmt.Sprintf(
		query,
		constants.ACCOUNT_TABLE,
	)

	return db.Prepare(facQuery)
}

func initTransactionWithRetry(ctx context.Context, producer kev.KevProducer) error {
	var err error
	counter := 0
	for range 5 {
		err = producer.InitTransactions(ctx)
		if err == nil {
			return nil
		}
		time.Sleep(time.Second * 2) // atau exponential backoff bisa dipakai

		counter++
	}
	return fmt.Errorf("init transactions failed after %d attempts: %w", counter, err)
}

func NewAuthRepo(
	sri schrgs.SchemaRegisteryInstance,
	pgi pkg.PostgresInstance,
	avri avr.AvrSerdeInstance,
	kv kev.Kafka,
) (rp authRepo, _ error) {
	dbaddrs := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DBUSER"),
		os.Getenv("DBPASS"),
		os.Getenv("DBHOST"),
		os.Getenv("DBPORT"),
		os.Getenv("DBAUTH"),
	)

	db, err := pkg.OpenPostgres(dbaddrs, pgi)
	if err != nil {
		return rp, err
	}

	rp.db = db

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

	cfg := map[kev.ConfigKeyKafka]string{
		kev.BOOTSTRAP_SERVERS:  os.Getenv("KAKFKABROKER"),
		kev.ACKS:               "all",
		kev.ENABLE_IDEMPOTENCE: "true",
		kev.COMPRESION_TYPE:    "snappy",
		kev.RETRIES:            "5",
		kev.RETRY_BACKOFF_MS:   "100",
		kev.LINGER_MS:          "5",
		kev.MAX_IN_FLIGHT_REQUESTS_PER_CONNECTION: "5",
		kev.TRANSACTIONAL_ID:                      TRANSACTIONAL_ID,
	}

	srgs, err := schrgs.NewSchemaRegistery(os.Getenv("SCHEMAREGISTERYADDR"), sri)
	if err != nil {
		return rp, err
	}

	rp.schemaRegisteryClient = srgs.Client()

	serializer, err := avri.NewGenericSerializer(
		rp.schemaRegisteryClient.Client(),
		avr.ValueSerde,
		avr.NewSerializerConfig(),
	)
	if err != nil {
		return rp, err
	}

	rp.avroSerializer = serializer

	pool := kev.NewKafkaProducerPool(kv)

	producer, err := pool.Producer(cfg)
	if err != nil {
		return rp, err
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		15*time.Second,
	)

	defer cancel()

	if err := initTransactionWithRetry(ctx, producer); err != nil {
		return rp, err
	}

	rp.authProducer = producer

	return rp, nil
}

func (rp authRepo) CloseRepo() {
	rp.schemaRegisteryClient.Client().Close()
	rp.avroSerializer.Close()
	rp.authProducer.Close()
	rp.db.Close()
}

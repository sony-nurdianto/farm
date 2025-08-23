package repository

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/sony-nurdianto/farm/auth/internal/concurrent"
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
	schemaRegisteryClient schrgs.SchrgsClient
	db                    pkg.PostgresDatabase
	avroSerializer        avr.AvrSerializer
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

func send(
	ctx context.Context,
	send chan any,
	recv any,
) {
	select {
	case <-ctx.Done():
		return
	case send <- recv:
	}
}

type repoPgDb struct {
	db                 pkg.PostgresDatabase
	getUserByEmailStmt pkg.Stmt
}

func initPostgresDb(ctx context.Context, pgi pkg.PostgresInstance) <-chan any {
	out := make(chan any, 1)

	go func() {
		defer close(out)
		var res concurrent.Result[pkg.PostgresDatabase]

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
			res.Error = err
			send(ctx, out, res)
			return
		}

		res.Value = db
		send(ctx, out, res)
	}()

	return out
}

func prepareDb(ctx context.Context, dbChan <-chan any) <-chan any {
	out := make(chan any, 1)
	go func() {
		defer close(out)
		var res concurrent.Result[repoPgDb]

		dbcv, ok := <-dbChan
		if !ok {
			return
		}

		dbres, ok := dbcv.(concurrent.Result[pkg.PostgresDatabase])
		if !ok {
			res.Error = errors.New("Wrong type data")
			send(ctx, out, res)
			return
		}

		if dbres.Error != nil {
			res.Error = dbres.Error
			send(ctx, out, res)
			return
		}

		res.Value.db = dbres.Value

		gue, err := prepareStmt(constants.QUERY_GET_USER_BY_EMAIL, dbres.Value)
		if err != nil {
			res.Error = err
			send(ctx, out, res)
			return
		}

		res.Value.getUserByEmailStmt = gue

		send(ctx, out, res)
	}()
	return out
}

func initSchemaRegistery(ctx context.Context, sri schrgs.SchemaRegisteryInstance) <-chan any {
	out := make(chan any, 1)
	go func() {
		defer close(out)
		var res concurrent.Result[schrgs.SchrgsClient]

		client, err := sri.NewClient(
			sri.NewConfig(
				os.Getenv("SCHEMAREGISTERYADDR"),
			),
		)
		if err != nil {
			res.Error = err
			send(ctx, out, res)
			return
		}

		res.Value = client
		send(ctx, out, res)
	}()
	return out
}

type schemaRegistryPair struct {
	Serializer avr.AvrSerializer
	Client     schrgs.SchrgsClient
}

func schemaNSerializerPipe(ctx context.Context, avri avr.AvrSerdeInstance, clientChan <-chan any) <-chan any {
	out := make(chan any, 1)
	go func() {
		defer close(out)
		var res concurrent.Result[schemaRegistryPair]

		client := <-clientChan

		schRes, ok := client.(concurrent.Result[schrgs.SchrgsClient])
		if !ok {
			res.Error = errors.New("Wrong type data")
			send(ctx, out, res)
			return
		}

		if schRes.Error != nil {
			res.Error = schRes.Error
			send(ctx, out, res)
			return
		}

		seri, err := avri.NewGenericSerializer(
			schRes.Value,
			avr.ValueSerde,
			avr.NewSerializerConfig(),
		)
		if err != nil {
			res.Error = err
			send(ctx, out, res)
			return
		}
		res.Value.Serializer = seri
		res.Value.Client = schRes.Value
		send(ctx, out, res)
	}()
	return out
}

func initTransactionWithRetry(ctx context.Context, producer kev.KevProducer) error {
	var err error
	counter := 0
	for range 5 {
		err = producer.InitTransactions(ctx)
		if err == nil {
			return nil
		}
		time.Sleep(time.Second * 2)

		counter++
	}
	return fmt.Errorf("init transactions failed after %d attempts: %w", counter, err)
}

func initAuthProducer(ctx context.Context, kv kev.Kafka) <-chan any {
	out := make(chan any, 1)
	go func() {
		defer close(out)
		var res concurrent.Result[kev.KevProducer]

		pool := kev.NewKafkaProducerPool(kv)
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

		producer, err := pool.Producer(cfg)
		if err != nil {
			res.Error = err
			send(ctx, out, res)
			return
		}

		ictx, cancel := context.WithTimeout(ctx, 15*time.Second)

		defer cancel()

		if err := initTransactionWithRetry(ictx, producer); err != nil {
			res.Error = err
			send(ctx, out, res)
			return
		}

		res.Value = producer
		send(ctx, out, res)
	}()
	return out
}

func NewAuthRepo(
	ctx context.Context,
	sri schrgs.SchemaRegisteryInstance,
	pgi pkg.PostgresInstance,
	avri avr.AvrSerdeInstance,
	kv kev.Kafka,
) (rp authRepo, _ error) {
	opsCtx, done := context.WithTimeout(ctx, time.Second*30)
	defer done()

	dbch := initPostgresDb(opsCtx, pgi)
	src := initSchemaRegistery(opsCtx, sri)

	chs := []<-chan any{
		prepareDb(opsCtx, dbch),
		schemaNSerializerPipe(opsCtx, avri, src),
		initAuthProducer(opsCtx, kv),
	}

	for v := range concurrent.FanIn(opsCtx, chs...) {
		switch res := v.(type) {
		case concurrent.Result[repoPgDb]:
			if res.Error != nil {
				return rp, res.Error
			}
			rp.db = res.Value.db
			rp.getUserByEmailStmt = res.Value.getUserByEmailStmt

		case concurrent.Result[schemaRegistryPair]:
			if res.Error != nil {
				return rp, res.Error
			}
			rp.schemaRegisteryClient = res.Value.Client
			rp.avroSerializer = res.Value.Serializer

		case concurrent.Result[kev.KevProducer]:
			if res.Error != nil {
				return rp, res.Error
			}
			rp.authProducer = res.Value
		}
	}

	return rp, nil
}

func (rp authRepo) CloseRepo() {
	rp.schemaRegisteryClient.Close()
	rp.avroSerializer.Close()
	rp.authProducer.Close()
	rp.db.Close()
}

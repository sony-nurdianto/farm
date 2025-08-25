package repo

import (
	"context"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/redis/go-redis/v9"
	"github.com/sony-nurdianto/farm/services/Events/farmer/constants"
	"github.com/sony-nurdianto/farm/services/Events/farmer/internal/models"
	"github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres/pkg"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/avr"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/kev"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/schrgs"
)

type Repo struct {
	schemaRegisteryClient schrgs.SchemaRegisteryClient
	avroDeserializer      avr.AvrDeserializer
	authConsumer          kev.KevConsumer
	authDB                authDB
	rdb                   *redis.Client
}

type authDB struct {
	db              pkg.PostgresDatabase
	updateEmailStmt pkg.Stmt
}

func NewAuthRepo(
	sri schrgs.SchemaRegisteryInstance,
	avri avr.AvrSerdeInstance,
	kv kev.Kafka,
	kevcfg map[kev.ConfigKeyKafka]string,
	pgi pkg.PostgresInstance,
	rdb *redis.Client,
) (ap Repo, _ error) {
	srgs, err := schrgs.NewSchemaRegistery(os.Getenv("SCHEMAREGISTERYADDR"), sri)
	if err != nil {
		return ap, err
	}
	ap.schemaRegisteryClient = srgs.Client()

	dese, err := avri.NewGenericDeserializer(
		ap.schemaRegisteryClient.Client(),
		serde.ValueSerde,
		avr.NewDeserializerConfig(),
	)
	if err != nil {
		return ap, err
	}

	ap.avroDeserializer = dese
	pool := kev.NewKafkaConsumerPool(kv)
	consumer, err := pool.Consumer(kevcfg)
	if err != nil {
		return ap, err
	}

	ap.authConsumer = consumer

	pgDB, err := pkg.OpenPostgres(os.Getenv("AUTH_DATABASE_ADDR"), pgi)
	if err != nil {
		return ap, err
	}

	ues, err := pgDB.Prepare(constants.AuthUpdateEmailAccount)
	if err != nil {
		return ap, err
	}

	athDB := authDB{
		db:              pgDB,
		updateEmailStmt: ues,
	}

	ap.authDB = athDB

	ap.rdb = rdb

	return ap, nil
}

func (ar Repo) SyncAccountsEmail(ctx context.Context, id, email string) error {
	row := ar.authDB.updateEmailStmt.QueryRowContext(ctx, email, id)
	if err := row.Err(); err != nil {
		return err
	}

	return nil
}

func (ar Repo) UpsertFarmerCache(ctx context.Context, key string, farmer models.Farmer) error {
	hset := ar.rdb.HSet(ctx, key, farmer)
	if _, err := hset.Result(); err != nil {
		return err
	}
	return nil
}

func (ar Repo) DeleteFarmerCache(ctx context.Context, key string) error {
	del := ar.rdb.Del(ctx, key)
	if _, err := del.Result(); err != nil {
		return err
	}
	return nil
}

func (ar Repo) DeserializerFarmer(topic string, payload []byte) (f models.Farmer, _ error) {
	if err := ar.avroDeserializer.DeserializeInto(topic, payload, &f); err != nil {
		return f, err
	}

	return f, nil
}

func (ar Repo) Consumer() kev.KevConsumer {
	return ar.authConsumer
}

func (ar Repo) CloseRepo() {
	ar.schemaRegisteryClient.Client().Close()
	ar.avroDeserializer.Close()
	ar.authConsumer.Close()
	ar.rdb.Close()
}

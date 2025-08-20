package repo

import (
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/avr"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/kev"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/schrgs"
)

type authRepo struct {
	schemaRegisteryClient schrgs.SchemaRegisteryClient
	avro                  avr.AvrSerdeInstance
	authConsumer          kev.KevConsumer
	rdb                   *redis.Client
}

func NewAuthRepo(
	sri schrgs.SchemaRegisteryInstance,
	avr avr.AvrSerdeInstance,
	kv kev.Kafka,
	cfg map[kev.ConfigKeyKafka]string,
	rdb *redis.Client,
) (ap authRepo, _ error) {
	srgs, err := schrgs.NewSchemaRegistery(os.Getenv("SCHEMAREGISTERYADDR"), sri)
	if err != nil {
		return ap, err
	}

	ap.schemaRegisteryClient = srgs.Client()
	ap.avro = avr

	pool := kev.NewKafkaConsumerPool(kv)
	consumer, err := pool.Consumer(cfg)

	ap.authConsumer = consumer
	ap.rdb = rdb

	return ap, nil
}

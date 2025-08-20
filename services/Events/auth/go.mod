module github.com/sony-nurdianto/farm/services/Events/auth

go 1.24.6

require github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev v0.0.0

require (
	github.com/actgardner/gogen-avro/v10 v10.2.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/confluentinc/confluent-kafka-go/v2 v2.11.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/heetch/avro v0.4.5 // indirect
	github.com/redis/go-redis/v9 v9.12.1 // indirect
	golang.org/x/oauth2 v0.18.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)

replace github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev => ../../../shared_lib/Go/kafkaev

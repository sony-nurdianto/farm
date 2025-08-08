module github.com/sony-nurdianto/farm/auth

go 1.24.5

require (
	github.com/google/uuid v1.6.0
	github.com/stretchr/testify v1.10.0
	golang.org/x/crypto v0.40.0
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/actgardner/gogen-avro/v10 v10.2.1 // indirect
	github.com/confluentinc/confluent-kafka-go v1.9.2 // indirect
	github.com/confluentinc/confluent-kafka-go/v2 v2.11.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/heetch/avro v0.4.5 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	golang.org/x/oauth2 v0.30.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres v0.0.0
	github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev v0.0.0
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250528174236-200df99c418a // indirect
	google.golang.org/grpc v1.74.2
)

replace (
	github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres => ../../../shared_lib/Go/database/postgres
	github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev => ../../../shared_lib/Go/kafkaev
)

module github.com/sony-nurdianto/farm/services/Grpc/ranch

go 1.25.0

require (
	github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres v0.0.0
	google.golang.org/grpc v1.75.1
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250707201910-8d1bb00bc6a7 // indirect
)

replace (
	github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres => ../../../shared_lib/Go/database/postgres
	github.com/sony-nurdianto/farm/shared_lib/Go/database/redis => ../../../shared_lib/Go/database/redis
	github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev => ../../../shared_lib/Go/kafkaev
	github.com/sony-nurdianto/farm/shared_lib/Go/observability => ../../../shared_lib/Go/observability
)

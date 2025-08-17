module github.com/sony-nurdianto/farm/auth

go 1.24.6

require (
	github.com/aead/chacha20poly1305 v0.0.0-20170617001512-233f39982aeb
	github.com/confluentinc/confluent-kafka-go/v2 v2.11.0
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/o1egl/paseto v1.0.0
	github.com/stretchr/testify v1.10.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.62.0
	go.opentelemetry.io/otel v1.37.0
	golang.org/x/crypto v0.40.0
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/actgardner/gogen-avro/v10 v10.2.1 // indirect
	github.com/aead/chacha20 v0.0.0-20180709150244-8b13a72661da // indirect
	github.com/aead/poly1305 v0.0.0-20180717145839-3fee0db0b635 // indirect
	github.com/cenkalti/backoff/v5 v5.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.1 // indirect
	github.com/heetch/avro v0.4.5 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/runtime v0.62.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc v0.13.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.37.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.37.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.37.0 // indirect
	go.opentelemetry.io/otel/log v0.13.0 // indirect
	go.opentelemetry.io/otel/metric v1.37.0 // indirect
	go.opentelemetry.io/otel/sdk v1.37.0 // indirect
	go.opentelemetry.io/otel/sdk/log v0.13.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.37.0 // indirect
	go.opentelemetry.io/otel/trace v1.37.0 // indirect
	go.opentelemetry.io/proto/otlp v1.7.0 // indirect
	golang.org/x/oauth2 v0.30.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250603155806-513f23925822 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres v0.0.0
	github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev v0.0.0
	github.com/sony-nurdianto/farm/shared_lib/Go/observability v0.0.0
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250603155806-513f23925822 // indirect
	google.golang.org/grpc v1.74.2
)

replace (
	github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres => ../../../shared_lib/Go/database/postgres
	github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev => ../../../shared_lib/Go/kafkaev
	github.com/sony-nurdianto/farm/shared_lib/Go/observability => ../../../shared_lib/Go/observability
)

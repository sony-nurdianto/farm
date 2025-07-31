module github.com/sony-nurdianto/farm/auth

go 1.24.5

require google.golang.org/protobuf v1.36.6

require github.com/lib/pq v1.10.9 // indirect

require (
	github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres v0.0.0
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250528174236-200df99c418a // indirect
	google.golang.org/grpc v1.74.2
)

replace github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres => ../shared_lib/Go/database/postgres

compose_test:
		podman compose -f compose.test.yaml up -d --build

compose_test_down:	
		podman compose -f compose.test.yaml  down 

gen_auth_proto:	
	 cd proto && buf generate --template ./auth/buf.gen.yaml --path ./auth/v1/auth.proto

gen_farmer_proto:	
	 cd proto && buf generate --template ./farmer/buf.gen.yaml --path ./farmer/v1/farmer.proto

gen_farm_proto:	
	 cd proto && buf generate --template ./farm/buf.gen.yaml --path ./farm/v1/farm.proto

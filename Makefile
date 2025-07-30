gen_auth_proto:	
	 cd proto && buf generate --template ./auth/buf.gen.yaml --path ./auth/v1/auth.proto

package intercpth

import (
	"log"

	"github.com/sony-nurdianto/farm/auth/internal/pbgen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func InterceptTokenValidate(req any) error {
	if req == nil {
		log.Printf("[AuthService] Nil request payload for TokenValidate")
		return status.Error(codes.InvalidArgument, "Expected Request is not nil")
	}

	dataRequest, ok := req.(*pbgen.TokenValidateRequest)
	if !ok {
		log.Printf("[AuthService] Invalid request type for TokenValidate - got: %T", req)
		return status.Error(codes.InvalidArgument, "Expected Request have type TokenValidateRequest Proto")
	}

	if len(dataRequest.Token) == 0 {
		log.Printf("[AuthService] Invalid request type for TokenValidate - Token is empty - does not meet requirements")
		return status.Error(codes.InvalidArgument, "Expected Request - Token is not empty")
	}

	return nil
}

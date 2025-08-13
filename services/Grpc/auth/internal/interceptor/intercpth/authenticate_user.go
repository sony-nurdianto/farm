package intercpth

import (
	"log"

	"github.com/sony-nurdianto/farm/auth/internal/pbgen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func InterceptAuthenticateUser(req any) error {
	if req == nil {
		log.Printf("[AuthService] Nil request payload for AuthenticateUser")
		return status.Error(codes.InvalidArgument, "Expected Request is not nil")
	}

	dataRequest, ok := req.(*pbgen.AuthenticateUserRequest)
	if !ok {
		log.Printf("[AuthService] Invalid request type for AuthenticateUser - got: %T", req)
		return status.Error(codes.InvalidArgument, "Expected Request have type AuthenticateUserRequest Proto")

	}

	if len(dataRequest.GetEmail()) == 0 {
		log.Printf("[AuthService] Invalid request type for AuthenticateUser - Email is empty - does not requirements")
		return status.Error(codes.InvalidArgument, "Email must not be empty")
	}

	if len(dataRequest.GetPassword()) == 0 {
		log.Printf("[AuthService] Invalid request type for AuthenticateUser - Password is empty - does not requirements")
		return status.Error(codes.InvalidArgument, "Password must not be empty")
	}

	log.Printf("[AuthService] AuthenticateUser request - Email: %s", dataRequest.Email)
	return nil
}

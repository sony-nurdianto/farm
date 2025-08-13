package intercpth

import (
	"log"

	"github.com/sony-nurdianto/farm/auth/internal/pbgen"
	"github.com/sony-nurdianto/farm/auth/internal/validator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func InterceptRegisterUser(req any) error {
	if req == nil {
		log.Printf("[AuthService] Nil request payload for Register")
		return status.Error(codes.InvalidArgument, "Expected Request is not nil")

	}

	dataRequest, ok := req.(*pbgen.RegisterUserRequest)
	if !ok {
		log.Printf("[AuthService] Invalid request type for Register - got: %T", req)
		return status.Error(codes.InvalidArgument, "Expected Request have type RegisterRequest Proto")
	}

	if !validator.ValidateEmail(dataRequest.Email) {
		log.Printf("[AuthService] Invalid request type for Register - Email Invalid - %s", dataRequest.Email)
		return status.Error(codes.InvalidArgument, "Email is not valid")
	}

	if !validator.ValidatePhone(dataRequest.PhoneNumber) {
		log.Printf("[AuthService] Invalid request type for Register - Phone Number Invalid - %s", dataRequest.PhoneNumber)
		return status.Error(codes.InvalidArgument, "Phone number is not valid")
	}

	if !validator.ValidatePassword(dataRequest.Password) {
		log.Printf("[AuthService] Invalid request type for Register - Password Invalid - does not meet complexity requirements")
		return status.Error(codes.InvalidArgument, "Password must be at least 8 characters, include 1 uppercase letter, 1 number, and 1 special character")
	}

	log.Printf("[AuthService] Register request - Email: %s, Phone: %s", dataRequest.Email, dataRequest.PhoneNumber)
	return nil
}

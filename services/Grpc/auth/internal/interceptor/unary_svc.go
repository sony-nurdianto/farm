package interceptor

import (
	"context"
	"log"
	"time"

	"github.com/sony-nurdianto/farm/auth/internal/pbgen"
	"github.com/sony-nurdianto/farm/auth/internal/validator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func AuthServiceUnaryInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {
	start := time.Now()

	log.Printf("[AuthService] Incoming request - Method: %s", info.FullMethod)
	switch info.FullMethod {
	case pbgen.AuthService_RegisterUser_FullMethodName:
		if req == nil {
			log.Printf("[AuthService] Nil request payload for Register")
			return nil, status.Error(codes.InvalidArgument, "Expected Request is not nil")

		}

		dataRequest, ok := req.(*pbgen.RegisterUserRequest)
		if !ok {
			log.Printf("[AuthService] Invalid request type for Register - got: %T", req)
			return nil, status.Error(codes.InvalidArgument, "Expected Request have type RegisterRequest Proto")
		}

		if !validator.ValidateEmail(dataRequest.Email) {
			log.Printf("[AuthService] Invalid request type for Register - Email Invalid - %s", dataRequest.Email)
			return nil, status.Error(codes.InvalidArgument, "Email is not valid")
		}

		if !validator.ValidatePhone(dataRequest.PhoneNumber) {
			log.Printf("[AuthService] Invalid request type for Register - Phone Number Invalid - %s", dataRequest.PhoneNumber)
			return nil, status.Error(codes.InvalidArgument, "Phone number is not valid")
		}

		if !validator.ValidatePassword(dataRequest.Password) {
			log.Printf("[AuthService] Invalid request type for Register - Password Invalid - does not meet complexity requirements")
			return nil, status.Error(codes.InvalidArgument, "Password must be at least 8 characters, include 1 uppercase letter, 1 number, and 1 special character")
		}

		log.Printf("[AuthService] Register request - Email: %s, Phone: %s", dataRequest.Email, dataRequest.PhoneNumber)
	case pbgen.AuthService_AuthenticateUser_FullMethodName:
		if req == nil {
			log.Printf("[AuthService] Nil request payload for AuthenticateUser")
			return nil, status.Error(codes.InvalidArgument, "Expected Request is not nil")
		}

		dataRequest, ok := req.(*pbgen.AuthenticateUserRequest)
		if !ok {
			log.Printf("[AuthService] Invalid request type for AuthenticateUser - got: %T", req)
			return nil, status.Error(codes.InvalidArgument, "Expected Request have type AuthenticateUserRequest Proto")

		}

		if len(dataRequest.GetEmail()) == 0 {
			log.Printf("[AuthService] Invalid request type for AuthenticateUser - Email is empty - does not requirements")
			return nil, status.Error(codes.InvalidArgument, "Email must not be empty")
		}

		if len(dataRequest.GetPassword()) == 0 {
			log.Printf("[AuthService] Invalid request type for AuthenticateUser - Password is empty - does not requirements")
			return nil, status.Error(codes.InvalidArgument, "Password must not be empty")
		}

		log.Printf("[AuthService] AuthenticateUser request - Email: %s", dataRequest.Email)
	}

	resp, err = handler(ctx, req)

	duration := time.Since(start)
	if err != nil {
		log.Printf("[AuthService] Method %s failed - Error: %v - Duration: %v", info.FullMethod, err, duration)
	} else {
		log.Printf("[AuthService] Method %s succeeded - Duration: %v", info.FullMethod, duration)
	}

	return resp, err
}

package interceptor

import (
	"context"
	"log"
	"time"

	"github.com/sony-nurdianto/farm/auth/internal/pbgen"
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
	case pbgen.AuthService_Register_FullMethodName:
		if req == nil {
			log.Printf("[AuthService] Nil request payload for Register")
			return nil, status.Error(codes.InvalidArgument, "Expected Request is not nil")

		}

		dataRequest, ok := req.(*pbgen.RegisterRequest)
		if !ok {
			log.Printf("[AuthService] Invalid request type for Register - got: %T", req)
			return nil, status.Error(codes.InvalidArgument, "Expected Request have type RegisterRequest Proto")
		}

		log.Printf("[AuthService] Register request - Email: %s, Phone: %s", dataRequest.Email, dataRequest.PhoneNumber)
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

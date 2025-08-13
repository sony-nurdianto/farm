package interceptor

import (
	"context"
	"log"
	"time"

	"github.com/sony-nurdianto/farm/auth/internal/interceptor/intercpth"
	"github.com/sony-nurdianto/farm/auth/internal/pbgen"
	"google.golang.org/grpc"
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
		if err := intercpth.InterceptRegisterUser(req); err != nil {
			return nil, err
		}

	case pbgen.AuthService_AuthenticateUser_FullMethodName:
		if err := intercpth.InterceptAuthenticateUser(req); err != nil {
			return nil, err
		}
	case pbgen.AuthService_TokenValidate_FullMethodName:
		if err := intercpth.InterceptTokenValidate(req); err != nil {
			return nil, err
		}
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

package unit_test

import (
	"context"
	"testing"

	"github.com/sony-nurdianto/farm/auth/internal/interceptor"
	"github.com/sony-nurdianto/farm/auth/internal/pbgen"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestUnaryInterceptorRequestIsNil(t *testing.T) {
	handler := func(ctx context.Context, req any) (any, error) {
		return &pbgen.RegisterResponse{
			Status: "Success",
			Msg:    "Success Register Response",
		}, nil
	}

	ctx := context.Background()

	info := &grpc.UnaryServerInfo{
		FullMethod: pbgen.AuthService_Register_FullMethodName,
	}

	_, err := interceptor.AuthServiceUnaryInterceptor(
		ctx,
		nil,
		info,
		handler,
	)

	assert.Error(t, err)
	assert.ErrorContains(t, err, "Expected Request is not nil")
}

func TestUnaryInterceptorRequestIsNotDefine(t *testing.T) {
	handler := func(ctx context.Context, req any) (any, error) {
		return &pbgen.RegisterResponse{
			Status: "Success",
			Msg:    "Success Register Response",
		}, nil
	}

	ctx := context.Background()
	req := &struct{}{}

	info := &grpc.UnaryServerInfo{
		FullMethod: pbgen.AuthService_Register_FullMethodName,
	}

	_, err := interceptor.AuthServiceUnaryInterceptor(
		ctx,
		req,
		info,
		handler,
	)

	assert.Error(t, err)
	assert.ErrorContains(t, err, "Expected Request have type RegisterRequest Proto")
}

func TestUnaryInterceptor(t *testing.T) {
	handler := func(ctx context.Context, req any) (any, error) {
		return &pbgen.RegisterResponse{
			Status: "Success",
			Msg:    "Success Register Response",
		}, nil
	}

	ctx := context.Background()
	req := &pbgen.RegisterRequest{
		FullName:    "Sony",
		Email:       "Sony@gmail.com",
		PhoneNumber: "+62851588206",
		Password:    "SomePassword",
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: pbgen.AuthService_Register_FullMethodName,
	}

	resp, err := interceptor.AuthServiceUnaryInterceptor(
		ctx,
		req,
		info,
		handler,
	)

	registerResp, ok := resp.(*pbgen.RegisterResponse)
	assert.True(t, ok)
	assert.NoError(t, err)

	assert.Equal(t, registerResp.Msg, "Success Register Response")
	assert.Equal(t, registerResp.Status, "Success")
}

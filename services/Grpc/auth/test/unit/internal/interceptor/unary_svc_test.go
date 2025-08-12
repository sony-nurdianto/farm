package unit_test

import (
	"context"
	"testing"

	"github.com/sony-nurdianto/farm/auth/internal/interceptor"
	"github.com/sony-nurdianto/farm/auth/internal/pbgen"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUnaryInterceptorRequestIsNil(t *testing.T) {
	handler := func(ctx context.Context, req any) (any, error) {
		return &pbgen.RegisterUserResponse{
			Status: "Success",
			Msg:    "Success Register Response",
		}, nil
	}

	ctx := context.Background()

	info := &grpc.UnaryServerInfo{
		FullMethod: pbgen.AuthService_RegisterUser_FullMethodName,
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
		return &pbgen.RegisterUserResponse{
			Status: "Success",
			Msg:    "Success Register Response",
		}, nil
	}

	ctx := context.Background()
	req := &struct{}{}

	info := &grpc.UnaryServerInfo{
		FullMethod: pbgen.AuthService_RegisterUser_FullMethodName,
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

func TestAuthServiceUnaryInterceptor(t *testing.T) {
	handler := func(ctx context.Context, req any) (any, error) {
		return &pbgen.RegisterUserResponse{
			Status: "Success",
			Msg:    "Success Register Response",
		}, nil
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: pbgen.AuthService_RegisterUser_FullMethodName,
	}

	t.Run("Success", func(t *testing.T) {
		req := &pbgen.RegisterUserRequest{
			FullName:    "Sony",
			Email:       "Sony@gmail.com",
			PhoneNumber: "+62851588206",
			Password:    "Some@P4assword",
		}

		resp, err := interceptor.AuthServiceUnaryInterceptor(
			context.Background(),
			req,
			info,
			handler,
		)

		registerResp, ok := resp.(*pbgen.RegisterUserResponse)
		assert.True(t, ok)
		assert.NoError(t, err)
		assert.Equal(t, "Success", registerResp.Status)
		assert.Equal(t, "Success Register Response", registerResp.Msg)
	})

	t.Run("Invalid Email", func(t *testing.T) {
		req := &pbgen.RegisterUserRequest{
			FullName:    "Sony",
			Email:       "invalid-email",
			PhoneNumber: "+62851588206",
			Password:    "Some@P4assword",
		}

		resp, err := interceptor.AuthServiceUnaryInterceptor(
			context.Background(),
			req,
			info,
			handler,
		)

		assert.Nil(t, resp)
		assert.Error(t, err)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("Invalid Phone", func(t *testing.T) {
		req := &pbgen.RegisterUserRequest{
			FullName:    "Sony",
			Email:       "sony@gmail.com",
			PhoneNumber: "invalid-phone",
			Password:    "Some@P4assword",
		}

		resp, err := interceptor.AuthServiceUnaryInterceptor(
			context.Background(),
			req,
			info,
			handler,
		)

		assert.Nil(t, resp)
		assert.Error(t, err)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("Invalid Password", func(t *testing.T) {
		req := &pbgen.RegisterUserRequest{
			FullName:    "Sony",
			Email:       "sony@gmail.com",
			PhoneNumber: "+62851588206",
			Password:    "nopass", // invalid: too short, no uppercase, no special char
		}

		resp, err := interceptor.AuthServiceUnaryInterceptor(
			context.Background(),
			req,
			info,
			handler,
		)

		assert.Nil(t, resp)
		assert.Error(t, err)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})
}

func TestUnaryInterceptor(t *testing.T) {
	handler := func(ctx context.Context, req any) (any, error) {
		return &pbgen.RegisterUserResponse{
			Status: "Success",
			Msg:    "Success Register Response",
		}, nil
	}

	ctx := context.Background()
	req := &pbgen.RegisterUserRequest{
		FullName:    "Sony",
		Email:       "Sony@gmail.com",
		PhoneNumber: "+62851588206",
		Password:    "Some@P4assword",
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: pbgen.AuthService_RegisterUser_FullMethodName,
	}

	resp, err := interceptor.AuthServiceUnaryInterceptor(
		ctx,
		req,
		info,
		handler,
	)

	registerResp, ok := resp.(*pbgen.RegisterUserResponse)
	assert.True(t, ok)
	assert.NoError(t, err)

	assert.Equal(t, registerResp.Msg, "Success Register Response")
	assert.Equal(t, registerResp.Status, "Success")
}

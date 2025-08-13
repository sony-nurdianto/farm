package unit_test

import (
	"context"
	"testing"
	"time"

	"github.com/sony-nurdianto/farm/auth/internal/interceptor"
	"github.com/sony-nurdianto/farm/auth/internal/pbgen"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestUnaryInterceptorRegisterUserRequestIsNil(t *testing.T) {
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

func TestUnaryInterceptorRegisterUserRequestIsNotDefine(t *testing.T) {
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

func TestAuthServiceUnaryInterceptorRegisterUserValidateRequest(t *testing.T) {
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

func TestUnaryInterceptorRegisterUserSucesss(t *testing.T) {
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

func TestUnaryInterceptorAuthenticateUserErrorRequestNil(t *testing.T) {
	handler := func(ctx context.Context, req any) (any, error) {
		return nil, nil
	}

	ctx := context.Background()

	info := &grpc.UnaryServerInfo{
		FullMethod: pbgen.AuthService_AuthenticateUser_FullMethodName,
	}

	resp, err := interceptor.AuthServiceUnaryInterceptor(
		ctx,
		nil,
		info,
		handler,
	)

	registerResp, ok := resp.(*pbgen.AuthenticateUserResponse)
	assert.False(t, ok)
	assert.Error(t, err)
	assert.Nil(t, registerResp)
}

func TestUnaryInterceptorAuthenticateUserErrorEmailEmpty(t *testing.T) {
	handler := func(ctx context.Context, req any) (any, error) {
		return nil, nil
	}

	req := &pbgen.AuthenticateUserRequest{
		Email:    "Sony@gmail.com",
		Password: "",
	}

	ctx := context.Background()

	info := &grpc.UnaryServerInfo{
		FullMethod: pbgen.AuthService_AuthenticateUser_FullMethodName,
	}

	resp, err := interceptor.AuthServiceUnaryInterceptor(
		ctx,
		req,
		info,
		handler,
	)

	registerResp, ok := resp.(*pbgen.AuthenticateUserResponse)
	assert.False(t, ok)
	assert.Error(t, err)
	assert.Nil(t, registerResp)
}

func TestUnaryInterceptorAuthenticateUserErrorPasswordEmpty(t *testing.T) {
	handler := func(ctx context.Context, req any) (any, error) {
		return nil, nil
	}

	req := &pbgen.AuthenticateUserRequest{
		Email:    "",
		Password: "Some@P4assword",
	}

	ctx := context.Background()

	info := &grpc.UnaryServerInfo{
		FullMethod: pbgen.AuthService_AuthenticateUser_FullMethodName,
	}

	resp, err := interceptor.AuthServiceUnaryInterceptor(
		ctx,
		req,
		info,
		handler,
	)

	registerResp, ok := resp.(*pbgen.AuthenticateUserResponse)
	assert.False(t, ok)
	assert.Error(t, err)
	assert.Nil(t, registerResp)
}

func TestUnaryInterceptorAuthenticateUserSucesss(t *testing.T) {
	handler := func(ctx context.Context, req any) (any, error) {
		return &pbgen.AuthenticateUserResponse{
			Token:     "Token",
			Status:    "Success",
			Msg:       "Success AuthenticateUser",
			IssuedAt:  timestamppb.Now(),
			ExpiresAt: timestamppb.New(time.Now().Add(time.Hour * 1)),
		}, nil
	}

	ctx := context.Background()
	req := &pbgen.AuthenticateUserRequest{
		Email:    "Sony@gmail.com",
		Password: "Some@P4assword",
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: pbgen.AuthService_AuthenticateUser_FullMethodName,
	}

	resp, err := interceptor.AuthServiceUnaryInterceptor(
		ctx,
		req,
		info,
		handler,
	)

	registerResp, ok := resp.(*pbgen.AuthenticateUserResponse)
	assert.True(t, ok)
	assert.NoError(t, err)

	assert.Equal(t, registerResp.Msg, "Success AuthenticateUser")
	assert.Equal(t, registerResp.Status, "Success")
}

func TestUnaryInterceptorTokenValidate_ErrorCases(t *testing.T) {
	handler := func(ctx context.Context, req any) (any, error) {
		t.Fatal("Handler should not be called on invalid request")
		return nil, nil
	}

	ctx := context.Background()
	info := &grpc.UnaryServerInfo{
		FullMethod: pbgen.AuthService_TokenValidate_FullMethodName,
	}

	t.Run("Nil request", func(t *testing.T) {
		resp, err := interceptor.AuthServiceUnaryInterceptor(
			ctx,
			nil,
			info,
			handler,
		)
		assert.Nil(t, resp)
		assert.Error(t, err)
		st, _ := status.FromError(err)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("Wrong type request", func(t *testing.T) {
		resp, err := interceptor.AuthServiceUnaryInterceptor(
			ctx,
			"not a proto request",
			info,
			handler,
		)
		assert.Nil(t, resp)
		assert.Error(t, err)
		st, _ := status.FromError(err)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("Empty token", func(t *testing.T) {
		resp, err := interceptor.AuthServiceUnaryInterceptor(
			ctx,
			&pbgen.TokenValidateRequest{Token: ""},
			info,
			handler,
		)
		assert.Nil(t, resp)
		assert.Error(t, err)
		st, _ := status.FromError(err)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})
}

func TestUnaryInterceptorTokenValidateSucesss(t *testing.T) {
	subject := "subject"
	isuer := "isuer"
	experiedAt := timestamppb.New(time.Now().Add(time.Hour * 1))

	handler := func(ctx context.Context, req any) (any, error) {
		return &pbgen.TokenValidateResponse{
			Valid:     true,
			Msg:       "Token Is Valid",
			Subject:   &subject,
			ExpiresAt: experiedAt,
			Isuer:     &isuer,
		}, nil
	}

	ctx := context.Background()
	req := &pbgen.TokenValidateRequest{
		Token: "Token",
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: pbgen.AuthService_TokenValidate_FullMethodName,
	}

	resp, err := interceptor.AuthServiceUnaryInterceptor(
		ctx,
		req,
		info,
		handler,
	)

	res, ok := resp.(*pbgen.TokenValidateResponse)
	assert.True(t, ok)
	assert.NoError(t, err)
	assert.True(t, res.Valid)
	assert.Equal(t, res.Msg, "Token Is Valid")
	assert.Equal(t, *res.Isuer, isuer)
	assert.Equal(t, res.ExpiresAt, experiedAt)
	assert.Equal(t, *res.Subject, subject)
}

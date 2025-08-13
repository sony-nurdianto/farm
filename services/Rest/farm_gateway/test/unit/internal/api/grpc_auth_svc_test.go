package uni_test

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/api"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/pbgen"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/test/mocks"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGrpcAuthService_AuthUserRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("AuthUserRegister Success", func(t *testing.T) {
		mockAuthSvcClient := mocks.NewMockAuthServiceClient(ctrl)

		mockAuthSvcClient.EXPECT().
			RegisterUser(gomock.Any(), gomock.Any()).
			Return(&pbgen.RegisterUserResponse{
				Msg:    "Success",
				Status: "Success Register User",
			}, nil)

		req := &pbgen.RegisterUserRequest{}

		svc := api.NewGrpcService(mockAuthSvcClient)
		res, err := svc.AuthUserRegister(req)
		assert.NoError(t, err)
		assert.Equal(t, res.Status, "Success Register User")
		assert.Equal(t, res.Msg, "Success")
	})

	t.Run("AuthUserRegister Error", func(t *testing.T) {
		mockAuthSvcClient := mocks.NewMockAuthServiceClient(ctrl)

		mockAuthSvcClient.EXPECT().
			RegisterUser(gomock.Any(), gomock.Any()).
			Return(nil, errors.New("Error Request Register User"))

		req := &pbgen.RegisterUserRequest{}

		svc := api.NewGrpcService(mockAuthSvcClient)
		res, err := svc.AuthUserRegister(req)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestAuthUserSignIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("AuthUserSignIn Success", func(t *testing.T) {
		mockAuthSvcClient := mocks.NewMockAuthServiceClient(ctrl)

		experiedAt := timestamppb.New(time.Now().Add(time.Hour * 1))
		issuedAt := timestamppb.Now()

		mockAuthSvcClient.EXPECT().
			AuthenticateUser(gomock.Any(), gomock.Any()).
			Return(
				&pbgen.AuthenticateUserResponse{
					Token:     "Token",
					ExpiresAt: experiedAt,
					IssuedAt:  issuedAt,
					Msg:       "Success AuthenticateUser",
					Status:    "Success",
				}, nil)

		req := &pbgen.AuthenticateUserRequest{
			Email:    "sony@gmail.com",
			Password: "secret",
		}

		svc := api.NewGrpcService(mockAuthSvcClient)
		res, err := svc.AuthUserSignIn(req)
		assert.NoError(t, err)
		assert.Equal(t, res.Token, "Token")
		assert.Equal(t, res.ExpiresAt, experiedAt)
		assert.Equal(t, res.IssuedAt, issuedAt)
		assert.Equal(t, res.Msg, "Success AuthenticateUser")
		assert.Equal(t, res.Status, "Success")
	})

	t.Run("AuthUserSignIn Error", func(t *testing.T) {
		mockAuthSvcClient := mocks.NewMockAuthServiceClient(ctrl)

		mockAuthSvcClient.EXPECT().
			AuthenticateUser(gomock.Any(), gomock.Any()).
			Return(nil, errors.New("Error Request Authenticate User"))

		req := &pbgen.AuthenticateUserRequest{}

		svc := api.NewGrpcService(mockAuthSvcClient)
		res, err := svc.AuthUserSignIn(req)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestAuthTokenValidate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("AuthUserSignIn Success", func(t *testing.T) {
		mockAuthSvcClient := mocks.NewMockAuthServiceClient(ctrl)

		experiedAt := timestamppb.New(time.Now().Add(time.Hour * 1))
		issuer := "service-atuh"
		subject := "id"

		mockAuthSvcClient.EXPECT().
			TokenValidate(gomock.Any(), gomock.Any()).
			Return(
				&pbgen.TokenValidateResponse{
					Valid:     true,
					ExpiresAt: experiedAt,
					Msg:       "Success AuthenticateUser",
					Isuer:     &issuer,
					Subject:   &subject,
				}, nil)

		req := &pbgen.TokenValidateRequest{
			Token: "Token",
		}

		svc := api.NewGrpcService(mockAuthSvcClient)
		res, err := svc.AuthTokenValidate(req)
		assert.NoError(t, err)
		assert.True(t, res.Valid)
		assert.Equal(t, res.ExpiresAt, experiedAt)
		assert.Equal(t, *res.Isuer, issuer)
		assert.Equal(t, *res.Subject, subject)
	})

	t.Run("AuthUserRegister Error", func(t *testing.T) {
		mockAuthSvcClient := mocks.NewMockAuthServiceClient(ctrl)

		mockAuthSvcClient.EXPECT().
			TokenValidate(gomock.Any(), gomock.Any()).
			Return(nil, errors.New("Error ValidateToken"))

		req := &pbgen.TokenValidateRequest{}

		svc := api.NewGrpcService(mockAuthSvcClient)
		res, err := svc.AuthTokenValidate(req)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

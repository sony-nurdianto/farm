package unit_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/handlers/authh"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/models"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/pbgen"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/test/mocks"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestSignInHandlerErrorBodyParser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := fiber.New()

	mockAuthSvc := mocks.NewMockGrpcAuthService(ctrl)

	authHandler := authh.NewAuthHandler(mockAuthSvc)

	app.Post("/auth/signin", authHandler.SignIn)

	user := models.UserSignIn{
		Email:    "sony@gmail.com",
		Password: "secreet",
	}

	jsonReq, err := json.Marshal(&user)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/auth/signin",
		bytes.NewReader(jsonReq),
	)

	res, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestSignInHandlerAuthUserSignInErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := fiber.New()

	mockAuthSvc := mocks.NewMockGrpcAuthService(ctrl)

	mockAuthSvc.EXPECT().
		AuthUserSignIn(gomock.Any()).
		Return(
			nil, errors.New("Error something"),
		)

	authHandler := authh.NewAuthHandler(mockAuthSvc)

	app.Post("/auth/signin", authHandler.SignIn)

	user := models.UserSignIn{
		Email:    "sony@gmail.com",
		Password: "secreet",
	}

	jsonReq, err := json.Marshal(&user)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/auth/signin",
		bytes.NewReader(jsonReq),
	)
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestSignInHandlerSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := fiber.New()

	mockAuthSvc := mocks.NewMockGrpcAuthService(ctrl)

	mockAuthSvc.EXPECT().
		AuthUserSignIn(gomock.Any()).
		Return(
			&pbgen.AuthenticateUserResponse{
				Token:     "Token",
				IssuedAt:  timestamppb.Now(),
				ExpiresAt: timestamppb.New(time.Now().Add(time.Hour * 1)),
				Msg:       "Success Login Welcome",
				Status:    "Success",
			}, nil)

	authHandler := authh.NewAuthHandler(mockAuthSvc)

	app.Post("/auth/signin", authHandler.SignIn)

	user := models.UserSignIn{
		Email:    "sony@gmail.com",
		Password: "secreet",
	}

	jsonReq, err := json.Marshal(&user)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/auth/signin",
		bytes.NewReader(jsonReq),
	)
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, res.StatusCode)

	body, _ := io.ReadAll(res.Body)
	assert.Contains(t, string(body), `"status":"Success"`)
}

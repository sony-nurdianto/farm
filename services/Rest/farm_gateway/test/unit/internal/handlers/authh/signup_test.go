package unit_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/handlers/authh"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/models"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/pbgen"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestSignupHandlerBodyParserError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := fiber.New()

	mockAuthSvc := mocks.NewMockGrpcAuthService(ctrl)
	authHandler := authh.NewAuthHandler(mockAuthSvc)

	app.Post("/auth/signup", authHandler.SignUp)

	user := models.UserRegister{
		FullName:    "sony",
		Email:       "sony@gmail.com",
		PhoneNumber: "+62851444777",
		Password:    "Secret@1",
	}

	jsonReq, err := json.Marshal(&user)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/auth/signup",
		bytes.NewReader(jsonReq),
	)

	res, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestSignupHandlerAuthRegisterErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := fiber.New()

	mockAuthSvc := mocks.NewMockGrpcAuthService(ctrl)

	mockAuthSvc.EXPECT().
		AuthUserRegister(gomock.Any()).
		Return(nil, errors.New("Something Wrong"))

	authHandler := authh.NewAuthHandler(mockAuthSvc)

	app.Post("/auth/signup", authHandler.SignUp)

	user := models.UserRegister{
		FullName:    "sony",
		Email:       "sony@gmail.com",
		PhoneNumber: "+62851444777",
		Password:    "Secret@1",
	}

	jsonReq, err := json.Marshal(&user)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/auth/signup",
		bytes.NewReader(jsonReq),
	)
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestSignUpHandlerSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := fiber.New()

	mockAuthSvc := mocks.NewMockGrpcAuthService(ctrl)

	mockAuthSvc.EXPECT().
		AuthUserRegister(gomock.Any()).
		Return(&pbgen.RegisterUserResponse{
			Status: "Success",
			Msg:    "Success Register User",
		}, nil)

	authHandler := authh.NewAuthHandler(mockAuthSvc)

	app.Post("/auth/signup", authHandler.SignUp)

	user := models.UserRegister{
		FullName:    "sony",
		Email:       "sony@gmail.com",
		PhoneNumber: "+62851444777",
		Password:    "Secret@1",
	}

	jsonReq, err := json.Marshal(&user)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/auth/signup",
		bytes.NewReader(jsonReq),
	)
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, res.StatusCode)

	body, _ := io.ReadAll(res.Body)
	assert.Contains(t, string(body), `"status":"Success"`)
}

package uni_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/pbgen"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/routes"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestRoutes_AuthSignup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock GrpcAuthService
	mockGrpcSvc := mocks.NewMockGrpcAuthService(ctrl)
	mockGrpcSvc.EXPECT().
		AuthUserRegister(gomock.Any()).
		Return(&pbgen.RegisterUserResponse{
			Status: "Success",
			Msg:    "User registered",
		}, nil)

	app := fiber.New()
	routes := routes.NewRoutes(app, mockGrpcSvc)
	routes.Build()

	// Request body
	jsonReq := `{
        "FullName": "Sony",
        "Email": "sony@gmail.com",
        "PhoneNumber": "+62851",
        "Password": "Secret@1"
    }`

	req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewReader([]byte(jsonReq)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), `"status":"Success"`)
	assert.Contains(t, string(body), `"msg":"User registered"`)
}

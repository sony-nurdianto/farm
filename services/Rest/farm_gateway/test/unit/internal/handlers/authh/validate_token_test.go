package unit_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/handlers/authh"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/pbgen"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/test/mocks"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestAuthTokenBaseValidate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthSvc := mocks.NewMockGrpcAuthService(ctrl)
	handler := authh.NewAuthHandler(mockAuthSvc)

	t.Run("missing header", func(t *testing.T) {
		app := fiber.New()
		app.Get("/", handler.AuthTokenBaseValidate)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res, _ := app.Test(req)
		assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
	})

	t.Run("invalid header format", func(t *testing.T) {
		app := fiber.New()
		app.Get("/", handler.AuthTokenBaseValidate)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "InvalidFormatToken")
		res, _ := app.Test(req)
		assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
	})

	t.Run("grpc error", func(t *testing.T) {
		app := fiber.New()
		mockAuthSvc.EXPECT().
			AuthTokenValidate(gomock.Any()).
			Return(nil, errors.New("grpc fail"))

		app.Get("/", handler.AuthTokenBaseValidate)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer validtoken")
		res, _ := app.Test(req)
		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
	})

	t.Run("token invalid", func(t *testing.T) {
		app := fiber.New()
		mockAuthSvc.EXPECT().
			AuthTokenValidate(gomock.Any()).
			Return(&pbgen.TokenValidateResponse{
				Valid: false,
				Msg:   "expired",
			}, nil)

		app.Get("/", handler.AuthTokenBaseValidate)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer sometoken")
		res, _ := app.Test(req)
		assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
	})

	t.Run("token valid", func(t *testing.T) {
		app := fiber.New()
		mockAuthSvc.EXPECT().
			AuthTokenValidate(gomock.Any()).
			Return(&pbgen.TokenValidateResponse{
				Valid:     true,
				Msg:       "ok",
				Subject:   ptrString("user123"),
				Isuer:     ptrString("authservice"),
				ExpiresAt: timestamppb.New(time.Now().Add(1 * time.Hour)),
			}, nil)

		// Route pakai Next() biar ada endpoint setelah middleware
		app.Get("/", handler.AuthTokenBaseValidate, func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer validtoken")
		res, _ := app.Test(req)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)
	})
}

func ptrString(s string) *string {
	return &s
}

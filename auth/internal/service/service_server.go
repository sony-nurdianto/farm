package service

import (
	"context"

	"github.com/sony-nurdianto/farm/auth/internal/pbgen"
)

type AuthServiceServer struct {
	pbgen.UnimplementedAuthServiceServer
}

func (ass *AuthServiceServer) Register(
	ctx context.Context,
	in *pbgen.RegisterRequest,
) (out *pbgen.RegisterResponse, _ error) {
	return out, nil
}

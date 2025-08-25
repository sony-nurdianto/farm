package usecase

import (
	"context"
	"time"

	"github.com/sony-nurdianto/farm/services/Grpc/farmer/internal/models"
	"github.com/sony-nurdianto/farm/services/Grpc/farmer/internal/pbgen"
	"github.com/sony-nurdianto/farm/services/Grpc/farmer/internal/repo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type FarmerUsecase interface {
	GetUserByID(ctx context.Context, id string) (*pbgen.FarmerProfileResponse, error)
	UpdateUser(ctx context.Context, user *models.UpdateUsers) (models.Users, error)
}

type farmerUsecase struct {
	repo repo.FarmerRepo
}

func NewFarmerUseCase(repo repo.FarmerRepo) farmerUsecase {
	return farmerUsecase{
		repo: repo,
	}
}

func (fu farmerUsecase) GetUserByID(ctx context.Context, id string) (*pbgen.FarmerProfileResponse, error) {
	uCtx, done := context.WithTimeout(ctx, time.Second*15)
	defer done()

	data, err := fu.repo.GetUsersByIDFromCache(uCtx, id)
	if err != nil {
		return nil, err
	}

	registerAt, err := time.Parse(time.RFC3339Nano, data.RegisteredAt)
	if err != nil {
		return nil, err
	}
	updatedAt, err := time.Parse(time.RFC3339Nano, data.UpdatedAt)
	if err != nil {
		return nil, err
	}

	farmer := &pbgen.FarmerProfileResponse{
		Farmer: &pbgen.Farmer{
			Id:           data.ID,
			FullName:     data.FullName,
			Email:        data.Email,
			Phone:        data.Phone,
			Verified:     data.Verified,
			RegisteredAt: timestamppb.New(registerAt),
			UpdatedAt:    timestamppb.New(updatedAt),
		},
	}

	return farmer, nil
}

func (fu farmerUsecase) UpdateUser(ctx context.Context, user *models.UpdateUsers) (res models.Users, _ error) {
	uCtx, done := context.WithTimeout(ctx, time.Second*15)
	defer done()

	res, err := fu.repo.UpdateUser(uCtx, user)
	if err != nil {
		return res, err
	}

	return res, nil
}

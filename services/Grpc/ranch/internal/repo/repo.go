package repo

import (
	"context"
	"log"

	"github.com/sony-nurdianto/farm/services/Grpc/ranch/internal/concurent"
	"github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres/pkg"
)

type RanchRepo interface{}

type ranchRepo struct {
	ranchDB ranchDatabase
}

func NewRanchRepo(
	ctx context.Context,
	pgi pkg.PostgresInstance,
) (RanchRepo, error) {
	dbCh := initPostgresDB(ctx, pgi)
	chs := []<-chan any{
		prepareRanchDB(ctx, dbCh),
	}

	var repo ranchRepo

	for v := range concurent.FanIn(ctx, chs...) {
		switch res := v.(type) {
		case concurent.ResultType[ranchDatabase]:
			if res.Error != nil {
				log.Println(res.Error)
				return nil, res.Error
			}

			repo.ranchDB = res.Value
		}
	}

	return repo, nil
}

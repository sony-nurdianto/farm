package repo

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/sony-nurdianto/farm/services/Grpc/farmer/internal/concurent"
	"github.com/sony-nurdianto/farm/services/Grpc/farmer/internal/constants"
	"github.com/sony-nurdianto/farm/services/Grpc/farmer/internal/models"
	"github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres/pkg"
	"github.com/sony-nurdianto/farm/shared_lib/Go/database/redis"
)

type FarmerRepo interface {
	GetUsersByIDFromCache(ctx context.Context, id string) (farmer models.Users, _ error)
	UpdateUser(ctx context.Context, users *models.UpdateUsers) (models.Users, error)
}

type farmerRepo struct {
	farmerCache redis.RedisClient
	farmerDB    farmerDB
}

type farmerDB struct {
	db             pkg.PostgresDatabase
	updateUserStmt pkg.Stmt
}

func send(
	ctx context.Context,
	send chan any,
	recv any,
) {
	select {
	case <-ctx.Done():
		return
	case send <- recv:
	}
}

func initPostgresDB(ctx context.Context, pgi pkg.PostgresInstance, addr string) <-chan any {
	out := make(chan any, 1)

	go func() {
		defer close(out)
		var res concurent.Result[pkg.PostgresDatabase]

		db, err := pkg.OpenPostgres(addr, pgi)
		if err != nil {
			res.Error = err
			send(ctx, out, res)
			return
		}

		res.Value = db
		send(ctx, out, res)
	}()

	return out
}

func prepareFarmerStmt(ctx context.Context, dbChan <-chan any) <-chan any {
	out := make(chan any, 1)

	go func() {
		defer close(out)

		var res concurent.Result[farmerDB]

		dbCv := <-dbChan

		dbres, ok := dbCv.(concurent.Result[pkg.PostgresDatabase])
		if !ok {
			res.Error = errors.New("wrong data type")
			send(ctx, out, res)
			return
		}

		if dbres.Error != nil {
			res.Error = dbres.Error
			send(ctx, out, res)
			return
		}

		// ues = updateUsersStmt
		uus, err := dbres.Value.Prepare(constants.UserQueryUpdate)
		if err != nil {
			res.Error = err
			send(ctx, out, res)
			return
		}

		res.Value = farmerDB{
			db:             dbres.Value,
			updateUserStmt: uus,
		}

		send(ctx, out, res)
	}()

	return out
}

func initRedisDatabae(ctx context.Context, rdi redis.RedisInstance) <-chan any {
	out := make(chan any, 1)
	go func() {
		defer close(out)
		var res concurent.Result[redis.RedisClient]

		rdb := redis.NewRedisDB(rdi)
		rdc, err := rdb.InitRedisClient(context.Background(), &redis.FailoverOptions{
			MasterName:    os.Getenv("FARMER_REDIS_MASTER_NAME"),
			SentinelAddrs: []string{os.Getenv("SENTINEL_FARMER_REDIS_ADDR")},
			Username:      os.Getenv("FARMER_REDIS_MASTER_USER_NAME"),
			Password:      os.Getenv("FARMER_REDIS_MASTER_PASSWORD"),
			DB:            0,
		})
		if err != nil {
			res.Error = err
			send(ctx, out, res)
			return
		}

		res.Value = rdc
		send(ctx, out, res)
	}()

	return out
}

func NewFarmerRepo(
	ctx context.Context,
	pgi pkg.PostgresInstance,
	rdi redis.RedisInstance,
) (fr farmerRepo, err error) {
	opsCtx, done := context.WithTimeout(ctx, time.Second*30)
	defer done()

	farmerDBCh := initPostgresDB(ctx, pgi, os.Getenv("FARMER_DATABASE_ADDR"))

	chs := []<-chan any{
		prepareFarmerStmt(ctx, farmerDBCh),
		initRedisDatabae(opsCtx, rdi),
	}

	for v := range concurent.FanIn(opsCtx, chs...) {
		switch res := v.(type) {
		case concurent.Result[redis.RedisClient]:
			if res.Error != nil {
				return fr, res.Error
			}
			fr.farmerCache = res.Value
		case concurent.Result[farmerDB]:
			if res.Error != nil {
				return fr, res.Error
			}

			fr.farmerDB = res.Value
		}
	}

	return fr, nil
}

func (fr farmerRepo) GetUsersByIDFromCache(ctx context.Context, id string) (user models.Users, _ error) {
	hkey := fmt.Sprintf("users:%s", id)

	err := fr.farmerCache.HGetAll(ctx, hkey).Scan(&user)
	if err != nil {
		return user, nil
	}
	if user == (models.Users{}) {
		return user, errors.New("user is not existed")
	}
	return user, nil
}

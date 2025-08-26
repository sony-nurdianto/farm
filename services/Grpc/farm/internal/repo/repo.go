package repo

import (
	"context"
	"fmt"
	"os"

	"github.com/sony-nurdianto/farm/services/Grpc/farm/internal/concurent"
	"github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres/pkg"
	"github.com/sony-nurdianto/farm/shared_lib/Go/database/redis"
)

type farmRepo struct {
	farmCache redis.RedisClient
	farmDB    farmDB
}

const (
	CreateFarmStmtType string = "CreateFarmStmtType"
	UpdateFarmStmtType string = "UpdateFarmStmtType"
	DeleteFarmStmtType string = "DeleteFarmStmtType"
)

type farmStmt struct {
	stmt     pkg.Stmt
	stmtType string
}

type farmDB struct {
	db             pkg.PostgresDatabase
	createFarmStmt pkg.Stmt
	updateFarmStmt pkg.Stmt
	deleteFarmStmt pkg.Stmt
}

func initPostgresDB(
	ctx context.Context,
	pgi pkg.PostgresInstance,
	addr string,
) <-chan any {
	out := make(chan any, 1)

	go func() {
		defer close(out)
		var res concurent.Result[pkg.PostgresDatabase]

		db, err := pkg.OpenPostgres(addr, pgi)
		if err != nil {
			res.Error = err
			concurent.SendResult(ctx, out, res)
			return
		}

		res.Value = db
		concurent.SendResult(ctx, out, res)
	}()

	return out
}

func prepareStmt(ctx context.Context, db pkg.PostgresDatabase, query string, stmtType string) <-chan any {
	out := make(chan any, 1)
	go func() {
		defer close(out)
		var res concurent.Result[farmStmt]
		stmt, err := db.Prepare(query)
		if err != nil {
			res.Error = err
			concurent.SendResult(ctx, out, res)
			return
		}

		outRes := farmStmt{
			stmt:     stmt,
			stmtType: stmtType,
		}

		res.Value = outRes
		concurent.SendResult(ctx, out, res)
	}()
	return out
}

func prepareFarmDB(ctx context.Context, dbChan <-chan any) <-chan any {
	out := make(chan any, 1)
	go func() {
		defer close(out)
		var res concurent.Result[farmDB]
		dbcv, ok := <-dbChan
		if !ok {
			fmt.Println("chanel is close")
			return
		}

		db, ok := dbcv.(concurent.Result[pkg.PostgresDatabase])
		if !ok {
			res.Error = fmt.Errorf("expected dbChan have type pkg.PostgresDatabase but got %v", db)
			concurent.SendResult(ctx, out, res)
			return
		}

		if db.Error != nil {
			res.Error = db.Error
			concurent.SendResult(ctx, out, res)
			return
		}

		chs := []<-chan any{
			prepareStmt(ctx, db.Value, "createFarmStmt", CreateFarmStmtType),
			prepareStmt(ctx, db.Value, "updateFarmStmt", UpdateFarmStmtType),
			prepareStmt(ctx, db.Value, "deleteFarmStmt", DeleteFarmStmtType),
		}

		dbFarm := farmDB{
			db: db.Value,
		}

		for v := range concurent.FanIn(ctx, chs...) {
			vRes, ok := v.(concurent.Result[farmStmt])
			if !ok {
				res.Error = fmt.Errorf("expected receive value of type farmStmt but go %v", vRes)
				concurent.SendResult(ctx, out, res)
				return
			}

			if vRes.Error != nil {
				res.Error = vRes.Error
				concurent.SendResult(ctx, out, res)
				return
			}

			switch vRes.Value.stmtType {
			case CreateFarmStmtType:
				dbFarm.createFarmStmt = vRes.Value.stmt
			case UpdateFarmStmtType:
				dbFarm.updateFarmStmt = vRes.Value.stmt
			case DeleteFarmStmtType:
				dbFarm.deleteFarmStmt = vRes.Value.stmt
			}
		}

		res.Value = dbFarm
		concurent.SendResult(ctx, out, res)
	}()
	return out
}

func prepareFarmCache(ctx context.Context, rdi redis.RedisInstance) <-chan any {
	out := make(chan any, 1)
	go func() {
		defer close(out)

		var res concurent.Result[redis.RedisClient]

		rdb := redis.NewRedisDB(rdi)
		rdc, err := rdb.InitRedisClient(ctx, &redis.FailoverOptions{
			MasterName: os.Getenv("FARMER_REDIS_MASTER_NAME"),
			SentinelAddrs: []string{
				os.Getenv("SENTINEL_FARMER_REDIS_ADDR"),
				os.Getenv("SENTINEL_FARMER_REDIS_ADDR_2"),
			},
			Username: os.Getenv("FARMER_REDIS_MASTER_USER_NAME"),
			Password: os.Getenv("FARMER_REDIS_MASTER_PASSWORD"),
			DB:       0,
		})
		if err != nil {
			res.Error = err
			concurent.SendResult(ctx, out, res)
			return
		}

		res.Value = rdc
		concurent.SendResult(ctx, out, res)
	}()
	return out
}

func NewFarmRepo(
	ctx context.Context,
	pgi pkg.PostgresInstance,
	rdi redis.RedisInstance,
) (fr farmRepo, _ error) {
	dbCh := initPostgresDB(ctx, pgi, "addr")
	chs := []<-chan any{
		prepareFarmDB(ctx, dbCh),
		prepareFarmCache(ctx, rdi),
	}

	for v := range concurent.FanIn(ctx, chs...) {
		switch res := v.(type) {
		case concurent.Result[farmDB]:
			if res.Error != nil {
				return fr, res.Error
			}
			fr.farmDB = res.Value
		case concurent.Result[redis.RedisClient]:
			if res.Error != nil {
				return fr, res.Error
			}
			fr.farmCache = res.Value
		}
	}

	return fr, nil
}

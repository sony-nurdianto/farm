package repo

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/sony-nurdianto/farm/services/Grpc/ranch/internal/concurent"
	"github.com/sony-nurdianto/farm/services/Grpc/ranch/internal/constants"
	"github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres/pkg"
)

type RanchRepo interface{}

type ranchRepo struct {
	ranchDB ranchDatabase
}

type AnimalStmtType int

const (
	InsertAnimalStmtType AnimalStmtType = iota
	DeleteAnimalStmtType
	GetAnimalAscStmtType
	GetAnimalDescStmtType
	GetTotalAnimalStmtType
)

type animalStmt struct {
	stmtType AnimalStmtType
	stmt     pkg.Stmt
}

type ranchDatabase struct {
	db                 pkg.PostgresDatabase
	insertAnimalStmt   pkg.Stmt
	deleteAnimalStmt   pkg.Stmt
	getAnimalAscStmt   pkg.Stmt
	getAnimalDescStmt  pkg.Stmt
	getTotalAnimalStmt pkg.Stmt
}

func initPostgresDB(
	ctx context.Context,
	pgi pkg.PostgresInstance,
) <-chan any {
	out := make(chan any, 1)
	go func() {
		defer close(out)

		result := concurent.NewResult[pkg.PostgresDatabase]()

		db, err := pkg.OpenPostgres(os.Getenv("RANCH_DATABASE_ADDR"), pgi)
		if err != nil {
			result.Err(err).
				SendResult(ctx, out)
			return
		}

		result.Res(db).
			SendResult(ctx, out)
	}()
	return out
}

func prepareStmt(
	ctx context.Context,
	db pkg.PostgresDatabase,
	query string,
	stmType AnimalStmtType,
) <-chan any {
	out := make(chan any, 1)
	go func() {
		defer close(out)
		res := concurent.NewResult[animalStmt]()

		stmt, err := db.Prepare(query)
		if err != nil {
			res.Err(err).
				SendResult(ctx, out)
			return
		}

		stmtAnimal := animalStmt{
			stmtType: stmType,
			stmt:     stmt,
		}

		res.Res(stmtAnimal).
			SendResult(ctx, out)
	}()
	return out
}

func prepareRanchDB(
	ctx context.Context,
	chanDB <-chan any,
) <-chan any {
	out := make(chan any, 1)
	go func() {
		defer close(out)

		result := concurent.NewResult[ranchDatabase]()

		resDB, ok := <-chanDB
		if !ok {
			log.Println("Channel Db is close")
			return
		}

		db, ok := resDB.(concurent.ResultType[pkg.PostgresDatabase])
		if !ok {
			result.Err(fmt.Errorf("expected PostgresDatabase result, got %T", db)).
				SendResult(ctx, out)
			return
		}

		if db.Error != nil {
			result.Err(db.Error).
				SendResult(ctx, out)
			return
		}

		ranchDB := ranchDatabase{
			db: db.Value,
		}

		chs := []<-chan any{
			prepareStmt(ctx, ranchDB.db, constants.QueryInsertAnimal, InsertAnimalStmtType),
			prepareStmt(ctx, ranchDB.db, constants.QueryDeleteAnimal, DeleteAnimalStmtType),

			// GetAnimal
			prepareStmt(ctx, ranchDB.db, constants.QueryGetAnimalASC, GetAnimalAscStmtType),
			prepareStmt(ctx, ranchDB.db, constants.QueryGetAnimalDESC, GetAnimalDescStmtType),
			prepareStmt(ctx, ranchDB.db, constants.QueryTotalAnimals, GetTotalAnimalStmtType),
		}

		for v := range concurent.FanIn(ctx, chs...) {
			res, ok := v.(concurent.ResultType[animalStmt])
			if !ok {
				result.Err(fmt.Errorf("expected animalStmt result, got %T", res)).
					SendResult(ctx, out)
				return
			}

			if res.Error != nil {
				result.Err(res.Error).
					SendResult(ctx, out)
				return
			}

			switch res.Value.stmtType {
			case InsertAnimalStmtType:
				ranchDB.insertAnimalStmt = res.Value.stmt
			case DeleteAnimalStmtType:
				ranchDB.deleteAnimalStmt = res.Value.stmt
			case GetAnimalAscStmtType:
				ranchDB.getAnimalAscStmt = res.Value.stmt
			case GetAnimalDescStmtType:
				ranchDB.getAnimalDescStmt = res.Value.stmt
			case GetTotalAnimalStmtType:
				ranchDB.getTotalAnimalStmt = res.Value.stmt

			}

		}

		result.Res(ranchDB).
			SendResult(ctx, out)
	}()
	return out
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

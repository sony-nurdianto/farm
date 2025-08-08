package main

import (
	"log"

	"github.com/sony-nurdianto/farm/auth/internal/pbgen"
	"github.com/sony-nurdianto/farm/auth/internal/repository"
	"github.com/sony-nurdianto/farm/auth/internal/usecase"
	"github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres/pkg"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/avr"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/kev"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/schrgs"
)

func main() {
	pgi := pkg.NewPostgresInstance()
	rgs := schrgs.NewRegistery()
	avr := avr.NewAvrSerdeInstance()
	kv := kev.NewKafka()
	repo, err := repository.NewPostgresRepo(rgs, pgi, avr, kv)
	if err != nil {
		log.Fatalln(err)
	}

	uc := usecase.NewServiceUsecase(&repo)
	req := &pbgen.RegisterRequest{
		FullName:    "Sony Nurdianto",
		Email:       "sonynurdiant0@gmail.com",
		PhoneNumber: "+6285158820461",
		Password:    "Kohend@789",
	}
	res, err := uc.UserRegister(req)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(res)
}

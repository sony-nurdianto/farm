package main

import (
	"log"

	"github.com/sony-nurdianto/farm/auth/internal/pbgen"
	"github.com/sony-nurdianto/farm/auth/internal/repository"
	"github.com/sony-nurdianto/farm/auth/internal/usecase"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/schrgs"
)

func main() {
	registery := schrgs.NewRegistery()
	repo, err := repository.NewPostgresRepo(registery)
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

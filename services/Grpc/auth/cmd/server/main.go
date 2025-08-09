package main

import (
	"crypto/rand"
	"log"

	"github.com/sony-nurdianto/farm/auth/internal/encryption/codec"
	"github.com/sony-nurdianto/farm/auth/internal/encryption/passencrypt"
	"github.com/sony-nurdianto/farm/auth/internal/pbgen"
	"github.com/sony-nurdianto/farm/auth/internal/repository"
	"github.com/sony-nurdianto/farm/auth/internal/usecase"
	"github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres/pkg"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/avr"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/kev"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/schrgs"
)

func main() {
	repo, err := repository.NewAuthRepo(
		schrgs.NewRegistery(),
		pkg.NewPostgresInstance(),
		avr.NewAvrSerdeInstance(),
		kev.NewKafka(),
	)
	if err != nil {
		log.Fatalln(err)
	}

	uc := usecase.NewServiceUsecase(
		&repo,
		passencrypt.NewPassEncrypt(
			rand.Reader,
			codec.NewBase64Encoder(),
		),
	)

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

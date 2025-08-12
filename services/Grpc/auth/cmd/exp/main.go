package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/sony-nurdianto/farm/auth/internal/encryption/token"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}

	pastok := token.NewPassetoToken()
	tokhen := token.NewTokhan(pastok)

	tok, err := tokhen.CreateWebToken("sony")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(tok)

	tokDec, err := tokhen.VerifyWebToken(tok)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(tokDec)
}

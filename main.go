package main

import (
	"crywealth/credential"
	"crywealth/internal/wealth/handler/cli"
	"crywealth/internal/wealth/repo"
	"crywealth/internal/wealth/service"
	"crywealth/pkg/binance"
	"crywealth/storage"
	"log"
)

func main() {

	// ignore missing env file err
	_ = credential.LoadENVFromFile()

	binanceCreds, err := credential.BinanceSecret()
	if err != nil {
		log.Fatal(err)
	}

	client := binance.NewClient("https://api.binance.com", binanceCreds.APIKey, binanceCreds.APISecret)

	dataStorage := storage.New()
	if err := dataStorage.CreateIfNotExists(); err != nil {
		log.Fatal("err on create if not exists: ", err)
	}

	wealthRepo := repo.New(dataStorage)

	wealthSVC := service.New(client, wealthRepo)
	cliHandler := cli.New(wealthSVC)

	err = cliHandler.Print()
	if err != nil {
		log.Fatal(err)
	}
}

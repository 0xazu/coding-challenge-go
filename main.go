package main

import (
	"coding-challenge-go/app"
	"coding-challenge-go/app/cache"
	"coding-challenge-go/app/web"
	"coding-challenge-go/app/config"
	"log"
)

func main() {
	appConfig, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal(app.CannotLoadConfigFile, err)
	}

	inMemoryTransactionService := cache.NewTransactionService(appConfig)
	transactionHandler := web.TransactionHandler{TransactionService: inMemoryTransactionService}

	transactionHandler.Handle()
}
package main

import (
	"github.com/joho/godotenv"
	"log"
	"worker/internal/handler"
	"worker/internal/infra"
	"worker/internal/service"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	amqpClient := infra.NewAMQPClient()
	crawlerService := service.NewCrawlerService()
	crawlerHandler := handler.NewCrawlerHandler(amqpClient, crawlerService)
	crawlerHandler.Process()
}

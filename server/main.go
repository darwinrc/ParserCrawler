package main

import (
	"context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"log"
	"net/http"
	"server/internal/handler"
	"server/internal/infra"
	"server/internal/repo"
	"server/internal/service"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	router := mux.NewRouter()

	redisClient := infra.NewRedisClient()
	amqpClient := infra.NewAMQPClient()

	crawlerRepo := repo.NewCrawlerRepository(redisClient)
	crawlerService := service.NewCrawlerService(crawlerRepo, amqpClient)
	crawlerHandler := handler.NewCrawlerHandler(crawlerService)
	crawlerHandler.Attach(router)

	wsHandler := handler.NewWsHandler(context.Background(), crawlerService)
	wsHandler.Attach(router)

	// Separate goroutine for consuming the message queue and writing the crawled urls to the websocket connection
	go wsHandler.ProcessCrawledUrls()

	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type"})

	err := http.ListenAndServe(":5000", handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(router))
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}

	log.Println("Parser web crawler server listening on port 5000")

}

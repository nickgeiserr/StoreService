package main

import (
	"ClubmineStoreService/handlers"
	"ClubmineStoreService/logger"
	"ClubmineStoreService/services"
	"ClubmineStoreService/stores"
	"go.uber.org/zap"
	"os"
)

func init() {
	//if err := godotenv.Load(); err != nil {
	//	log.Fatalf("Error loading .env file: %v", err)
	//}
}

func main() {

	e := handlers.Echo()

	s := stores.New()

	ss := services.New(s)
	h := handlers.New(ss)

	handlers.SetAPI(e, h)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "1515"
	}

	logger.Fatal("failed to start server", zap.Error(e.Start(":"+PORT)))
}

package database

import (
	"ClubmineStoreService/logger"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
	"os"
)

var DB *pgxpool.Pool
var DbUrl string

func init() {

	err := logger.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := godotenv.Load(); err != nil {
		logger.Error("Oops! Sum went wrong", zap.Error(err))
	}

	if os.Getenv("LOCAL") == "true" {
		DbUrl = os.Getenv("LOCAL_DB_URL")
	} else {
		DbUrl = os.Getenv("DB_URL")
	}

	pool, err := pgxpool.New(context.Background(), DbUrl)
	if err != nil {
		logger.Error("Failed to create Database Pool.")
		return
	}

	DB = pool
}

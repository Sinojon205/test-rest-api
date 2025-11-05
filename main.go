package main

import (
	"log/slog"
	"test-rest-api/config"
	"test-rest-api/pkg/handler"
	"test-rest-api/pkg/repository"
	"test-rest-api/pkg/service"

	"github.com/jmoiron/sqlx"
	"github.com/sytallax/prettylog"
)

func main() {
	prettyHandler := prettylog.NewHandler(&slog.HandlerOptions{
		Level:       slog.LevelInfo,
		AddSource:   false,
		ReplaceAttr: nil,
	})
	logger := slog.New(prettyHandler)
	db, err := sqlx.Connect("postgres", config.Get().DatabaseDSN)
	if err != nil {
		logger.Error("[ERROR] failed to connect to db: %v", err)
		return
	}

	defer db.Close()
	logger.Info("db connected successfully")
	srv := new(Server)
	repos := repository.NewRepo(db)
	services := service.NewService(logger, repos /*, rdb, botAPI*/)
	handlers := handler.NewHandler(services, logger)
	router := handlers.InitRoutes()
	if err := srv.Run(config.Get().Port, router); err != nil {
		logger.Error("error occured while running http server:", err.Error())
	}
}

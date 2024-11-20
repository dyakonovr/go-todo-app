package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todo-app/internal/config"
	delivery "todo-app/internal/delivery/http"
	"todo-app/internal/repository"
	"todo-app/internal/server"
	"todo-app/internal/service"
	"todo-app/pkg/logger"

	"github.com/jackc/pgx/v5"
)

func Run(configDir string, configName string) {
	logger.Init()

	cfg, err := config.Init(configDir, configName)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	// Connection to database
	conn, err := pgx.Connect(context.Background(), cfg.Postgres.URI)
	if err != nil {
		logger.Errorf("Unable to connect to database: %v\n", err)
		return
	}
	defer conn.Close(context.Background())

	// Create repos, services, handler
	repos := repository.NewRepositories(conn)
	services := service.NewService(service.Deps{Repos: repos})
	handlers := delivery.NewHandler(services)

	// HTTP Server
	server := server.CreateNewServer(cfg, handlers.Init())

	go func() {
		if err := server.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	logger.Info("Server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := server.Stop(ctx); err != nil {
		logger.Errorf("failed to stop server: %v", err)
	}
}


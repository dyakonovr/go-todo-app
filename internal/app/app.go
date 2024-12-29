package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todo-app/internal/api"
	todosapi "todo-app/internal/api/todosApi"
	"todo-app/internal/config"
	"todo-app/internal/repo/todosRepo"
	"todo-app/internal/server"
	todosusecase "todo-app/internal/usecase/todosUsecase"
	"todo-app/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Run(configDir string, configName string) {
	ctx := context.Background()
	logger.Init()

	cfg, err := config.Init(configDir, configName)
	if err != nil {
		logger.Errorf("Error while compiling config: %v", err.Error())
		return
	}

	// Connection to database
	pool, err := pgxpool.New(ctx, cfg.Postgres.URI)
	if err != nil {
		logger.Errorf("Unable to connect to database: %v\n", err)
		return
	}

	defer pool.Close()

	// Create repos, services, apis
	todosRepo := todosRepo.New(pool)
	todosUsecase := todosusecase.New(todosRepo)
	todosApi := todosapi.New(todosUsecase)

	// HTTP Server
	server := server.CreateNewServer(cfg, api.Init([]api.ApiInterface{todosApi}))

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

	ctx, shutdown := context.WithTimeout(ctx, timeout)
	defer shutdown()

	if err := server.Stop(ctx); err != nil {
		logger.Errorf("failed to stop server: %v", err)
	}
}

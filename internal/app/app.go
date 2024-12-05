package app

import (
	"context"
	api "github.com/wachrusz/Back-End-API/internal/http"
	mydb "github.com/wachrusz/Back-End-API/internal/mydatabase"
	"github.com/wachrusz/Back-End-API/internal/server"
	"github.com/wachrusz/Back-End-API/pkg/cache"
	logger "github.com/zhukovrost/cadv_logger"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/wachrusz/Back-End-API/internal/config"
	"github.com/wachrusz/Back-End-API/internal/http/obhttp"
	"github.com/wachrusz/Back-End-API/internal/http/v1"
	"github.com/wachrusz/Back-End-API/internal/repository"
	"github.com/wachrusz/Back-End-API/internal/service"
	"github.com/wachrusz/Back-End-API/pkg/rabbit"
)

func Run(cfg *config.Config) error {
	l := logger.New("standard", true)

	l.Info("Starting application...")

	l.Info("Connecting to the database...")
	db, err := mydb.Init(cfg.GetDBURL())
	if err != nil {
		return err
	}
	defer db.Close()

	mydb.SetDB(db) // TODO: Избавиться от этой хуйни окончательно!

	l.Info("Connecting to redis...")
	redis, err := cache.New(cfg.Redis)
	if err != nil {
		return err
	}
	defer redis.Close()

	l.Info("Connecting to RabbitMQ...")
	mailer, err := rabbit.New(cfg.Rabbit, l)
	if err != nil {
		return err
	}

	l.Info("Initializing services...")
	deps := service.Dependencies{
		Repo:                  db,
		Mailer:                mailer,
		Cache:                 redis,
		AccessTokenDurMinutes: cfg.AccessTokenLifetime,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		l.Info("Received termination signal")
		cancel()
	}()

	if err := worker(ctx, cfg, deps, l); err != nil {
		return err
	}

	l.Info("Application stopped successfully")
	return nil
}

func worker(ctx context.Context, cfg *config.Config, deps service.Dependencies, logger *zap.Logger) error {
	services, err := service.NewServices(deps)
	if err != nil {
		return err
	}

	logger.Info("Initializing models...")
	models := repository.New(deps.Repo)

	logger.Info("Initializing handlers...", zap.Int64("rate_limit_per_second", cfg.RateLimitPerSecond))
	handlerV1 := v1.NewHandler(services, logger, models, deps.Cache, cfg.RateLimitPerSecond)
	handlerOB := obhttp.NewHandler(services, logger)

	logger.Info("Initializing routers...")
	router, docRouter, err := api.InitRouters(handlerV1, handlerOB, logger)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", router)
	mux.Handle("/swagger/", docRouter)
	mux.Handle("/docs/", docRouter)

	services.Users.InitActiveUsers()
	go services.Currency.ScheduleCurrencyUpdates()

	srv := server.NewServer(mux, logger, cfg.Server)
	return srv.Run(ctx)
}

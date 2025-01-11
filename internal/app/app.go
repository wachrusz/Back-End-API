package app

import (
	api "github.com/wachrusz/Back-End-API/internal/http"
	mydb "github.com/wachrusz/Back-End-API/internal/mydatabase"
	"github.com/wachrusz/Back-End-API/pkg/cache"
	logger "github.com/zhukovrost/cadv_logger"
	"go.uber.org/zap"
	"net/http"

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
		AccessTokenDurMinutes: cfg.AccessTokenLifetime,
	}

	services, err := service.NewServices(deps)
	if err != nil {
		return err
	}

	l.Info("Initializing models...")
	models := repository.New(db)

	l.Info("Initializing handlers...", zap.Int64("rate_limit_per_second", cfg.RateLimitPerSecond))
	handlerV1 := v1.NewHandler(services, l, models, redis, cfg.RateLimitPerSecond)
	handlerOB := obhttp.NewHandler(services, l)

	l.Info("Initializing routers...")
	router, docRouter, err := api.InitRouters(handlerV1, handlerOB, l)
	if err != nil {
		return err
	}

	http.Handle("/", router)
	http.Handle("/swagger/", docRouter)
	http.Handle("/docs/", docRouter)

	go services.Currency.ScheduleCurrencyUpdates()

	l.Info("Serving...")
	//changed tls hosting now everything works
	err = http.ListenAndServeTLS(":8080", cfg.CrtPath, cfg.KeyPath, nil)
	if err != nil {
		return err
	}

	return nil
}

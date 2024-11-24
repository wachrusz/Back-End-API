package app

import (
	"github.com/wachrusz/Back-End-API/internal/repository"
	"net/http"

	api "github.com/wachrusz/Back-End-API/internal/http"
	mydb "github.com/wachrusz/Back-End-API/internal/mydatabase"
	logger "github.com/zhukovrost/cadv_logger"

	"github.com/wachrusz/Back-End-API/internal/config"
	"github.com/wachrusz/Back-End-API/internal/http/obhttp"
	"github.com/wachrusz/Back-End-API/internal/http/v1"
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

	l.Info("Connecting to RabbitMQ...")
	mailer, err := rabbit.New(cfg.Rabbit, l)
	if err != nil {
		return err
	}

	deps := service.Dependencies{
		Repo:                  db,
		Mailer:                mailer,
		AccessTokenDurMinutes: cfg.AccessTokenLifetime,
	}

	l.Info("Initializing services...")
	services, err := service.NewServices(deps)
	if err != nil {
		return err
	}

	l.Info("Initializing models...")
	models := repository.New(db)

	handlerV1 := v1.NewHandler(services, l, models)
	handlerOB := obhttp.NewHandler(services, l)

	l.Info("Initializing routers...")
	router, docRouter, errR := api.InitRouters(handlerV1, handlerOB, l)
	services.Users.InitActiveUsers()

	if errR != nil {
		return errR
	}

	http.Handle("/", handlerV1.ContentTypeMiddleware(router))
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

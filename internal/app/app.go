package app

import (
	"net/http"

	"github.com/wachrusz/Back-End-API/internal/config"
	api "github.com/wachrusz/Back-End-API/internal/http"
	v1 "github.com/wachrusz/Back-End-API/internal/http/v1"
	mydb "github.com/wachrusz/Back-End-API/internal/mydatabase"
	"github.com/wachrusz/Back-End-API/internal/service"
	logger "github.com/zhukovrost/cadv_logger"
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

	deps := service.Dependencies{
		Repo: db,
	}

	l.Info("Initializing services...")
	services, err := service.NewServices(deps)
	if err != nil {
		return err
	}

	handler := v1.NewHandler(services, l)

	l.Info("Initializing routers...")
	router, docRouter, errR := api.InitRouters(handler, l)
	services.Users.InitActiveUsers()

	if errR != nil {
		return errR
	}

	http.Handle("/", handler.ContentTypeMiddleware(router))
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

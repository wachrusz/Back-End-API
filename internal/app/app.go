package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/wachrusz/Back-End-API/internal/config"
	api "github.com/wachrusz/Back-End-API/internal/http"
	"github.com/wachrusz/Back-End-API/internal/http/obhttp"
	v1 "github.com/wachrusz/Back-End-API/internal/http/v1"
	mydb "github.com/wachrusz/Back-End-API/internal/mydatabase"
	"github.com/wachrusz/Back-End-API/internal/repository"
	"github.com/wachrusz/Back-End-API/internal/server"
	"github.com/wachrusz/Back-End-API/internal/service"
	"github.com/wachrusz/Back-End-API/pkg/cache"
	"github.com/wachrusz/Back-End-API/pkg/rabbit"
	logger "github.com/zhukovrost/cadv_logger"
)

const (
	rateLimitKey = "rate_limit:total"
	stableChecks = 5
)

var (
	workerCount int64
	workerMu    sync.Mutex
	workerMap   = make(map[int64]context.CancelFunc)
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
	client, err := cache.New(cfg.Redis)
	if err != nil {
		return err
	}
	defer client.Close()

	l.Info("Connecting to RabbitMQ...")
	mailer, err := rabbit.New(cfg.Rabbit, l)
	if err != nil {
		return err
	}

	l.Info("Initializing services...")
	deps := service.Dependencies{
		Repo:                  db,
		Mailer:                mailer,
		Cache:                 client,
		AccessTokenDurMinutes: cfg.AccessTokenLifetime,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go monitorRateLimit(ctx, *cfg, deps, l)

	// Goroutine для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		l.Info("Received shutdown signal")
		cancel()
	}()

	if err := worker(ctx, cfg, deps, l); err != nil {
		return err
	}

	l.Info("Application stopped successfully")
	return nil
}

func worker(ctx context.Context, cfg *config.Config, deps service.Dependencies, logger *zap.Logger) error {
	logger.Info("Starting worker...")
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

func monitorRateLimit(ctx context.Context, cfg config.Config, deps service.Dependencies, logger *zap.Logger) {
	logger.Info("Started monitorRateLimit")
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	lowRateCount := 0

	for {
		logger.Info("Entrance")
		select {
		case <-ctx.Done():
			logger.Info("case Done")
			return
		case <-ticker.C:
			logger.Info("case <-ticker.C")
			value, err := deps.Cache.Get(ctx, rateLimitKey).Int64()
			if err != nil && !errors.Is(err, redis.Nil) {
				logger.Error("Failed to fetch rate limit value", zap.Error(err))
				continue
			}

			logger.Info(string(value))

			logger.Info(string(cfg.Workers.NewWorkerRPS * (workerCount + 1)))
			if value > cfg.Workers.NewWorkerRPS*(workerCount+1) {
				logger.Info("if value > cfg.Workers.NewWorkerRPS*(workerCount+1)")
				launchWorker(ctx, cfg, deps, logger)
				lowRateCount = 0
			} else if cfg.Workers.NewWorkerRPS*(workerCount+1) >= value && value >= cfg.Workers.NewWorkerRPS*workerCount {
				logger.Info("if cfg.Workers.NewWorkerRPS*(workerCount+1) >= value && value >= cfg.Workers.NewWorkerRPS*workerCount")
				lowRateCount = 0
			} else {
				logger.Info("else")
				lowRateCount++
				if lowRateCount >= stableChecks {
					terminateWorkers(context.Background(), cfg, logger)
				}
			}
		}
		logger.Info("lowRateCount: " + string(lowRateCount))
	}
}

func launchWorker(ctx context.Context, cfg config.Config, deps service.Dependencies, logger *zap.Logger) {
	logger.Info("Started launchWorker")
	workerMu.Lock()
	defer workerMu.Unlock()

	if workerCount >= cfg.Workers.MaxWorkers {
		logger.Warn("Max workers limit reached")
		return
	}

	workerCount++
	port := cfg.Server.Port + workerCount
	cfg.Server.Port = port

	logger.Info("workerCount, port" + string(workerCount) + string(port))

	// Запускаем воркера с новым контекстом
	workerCtx, workerCancel := context.WithCancel(ctx)

	go func() {
		defer workerCancel()
		err := worker(workerCtx, &cfg, deps, logger)
		if err != nil {
			logger.Error("Worker stopped with error", zap.Error(err))
		} else {
			logger.Info("Worker stopped gracefully", zap.Int64("port", cfg.Server.Port))
		}
	}()

	workerMap[port] = workerCancel
}

// Terminate all worker goroutines
func terminateWorkers(ctx context.Context, cfg config.Config, logger *zap.Logger) {
	workerMu.Lock()
	defer workerMu.Unlock()

	port := cfg.Server.Port + workerCount
	workerCount--

	cancel := workerMap[port]
	delete(workerMap, port)

	cancel()
	// Cancel each worker context and remove from channel
	logger.Info("Terminating worker", zap.Int64("port", port))

}

package server

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Config struct {
	Host          string        `yaml:"host"`
	Port          int           `yaml:"port"`
	IdleTimeout   time.Duration `yaml:"idle_timeout"`
	ReadTimeout   time.Duration `yaml:"read_timeout"`
	WriteTimeout  time.Duration `yaml:"write_timeout"`
	ShutdownGrace time.Duration `yaml:"shutdown_grace"`

	CrtPath string `yaml:"crt_path"`
	KeyPath string `yaml:"key_path"`
}

type Server struct {
	httpServer *http.Server
	logger     *zap.Logger
	cfg        Config
}

func NewServer(handler http.Handler, logger *zap.Logger, cfg Config) *Server {
	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:      handler,
		IdleTimeout:  cfg.IdleTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	return &Server{
		httpServer: httpServer,
		logger:     logger,
		cfg:        cfg,
	}
}

// start runs the server and listens for incoming requests.
func (s *Server) start() error {
	s.logger.Info("Starting server...", zap.Int("port", s.cfg.Port))
	return s.httpServer.ListenAndServeTLS(s.cfg.CrtPath, s.cfg.KeyPath)
}

// Run handles starting the server and managing graceful shutdown logic.
func (s *Server) Run(ctx context.Context) error {
	// Run server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- s.start()
	}()

	select {
	case err := <-errChan:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server error: %w", err)
		}
	case <-ctx.Done():
		// Trigger shutdown
		ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownGrace)
		defer cancel()
		if err := s.shutdown(ctx); err != nil {
			return fmt.Errorf("graceful shutdown failed: %w", err)
		}
		s.logger.Info("Server stopped gracefully", zap.Int("port", s.cfg.Port))
	}

	return nil
}

// shutdown gracefully shuts down the server with the given context.
func (s *Server) shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down server...", zap.Int("port", s.cfg.Port))
	return s.httpServer.Shutdown(ctx)
}

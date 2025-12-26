package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/RofaBR/Go-Usof/internal/config"
	"github.com/RofaBR/Go-Usof/internal/handler"
	"github.com/RofaBR/Go-Usof/internal/repositories"
	"github.com/RofaBR/Go-Usof/internal/router"
	"github.com/RofaBR/Go-Usof/internal/services"
	"github.com/RofaBR/Go-Usof/internal/storage/postgres"
	"github.com/RofaBR/Go-Usof/internal/storage/redis"
	"github.com/RofaBR/Go-Usof/pkg/logger"
	"github.com/gin-gonic/gin"
)

type App struct {
	config *config.Config
	logger *logger.Logger
	router *gin.Engine
	server *http.Server
	db     *postgres.Postgres
	redis  *redis.Redis
}

func New() (*App, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	log := logger.New(cfg.LogLevel)
	log.Info("initializing application")

	log.Info("connecting to PostgreSQL database")
	db, err := postgres.Run(context.Background(), cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Info("connecting to Redis")
	redisClient, err := redis.NewRedis(context.Background(), cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Info("Initializing repositories")
	repos := repositories.NewRepository(db, redisClient)

	log.Info("Initializing services")
	svc := services.NewServices(log, repos, cfg)

	log.Info("initializing handlers")
	handlers := handler.NewHandler(log, svc)

	gin.SetMode(cfg.Mode)
	r := router.SetupRouter(log, handlers)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	return &App{
		config: cfg,
		logger: log,
		router: r,
		server: server,
		db:     db,
		redis:  redisClient,
	}, nil
}

func (a *App) Run() error {
	a.logger.Info("starting server",
		"port", a.config.Port,
		"mode", a.config.Mode,
	)

	if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	a.logger.Info("shutting down server gracefully")

	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}

	a.logger.Info("closing database connection")
	a.db.Close()

	a.logger.Info("closing redis connection")
	if err := a.redis.Close(); err != nil {
		a.logger.Error("failed to close redis connection", "error", err)
	}

	a.logger.Info("server stopped successfully")
	return nil
}

func (a *App) Config() *config.Config {
	return a.config
}

func (a *App) Logger() *logger.Logger {
	return a.logger
}

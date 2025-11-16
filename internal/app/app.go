package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/leonidlivshits/Avito-MVP/config"
	"github.com/leonidlivshits/Avito-MVP/internal/log"
	"github.com/leonidlivshits/Avito-MVP/internal/infra/db"
	pgres "github.com/leonidlivshits/Avito-MVP/internal/infra/postgres"
	"github.com/leonidlivshits/Avito-MVP/internal/domain/service"
	"github.com/leonidlivshits/Avito-MVP/internal/usecase"
	httpapi "github.com/leonidlivshits/Avito-MVP/internal/api/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	cfg    *config.Config
	logger *log.Logger
	pool   *pgxpool.Pool
	server *http.Server
}

func NewApp(cfg *config.Config) (*App, error) {
	logger := log.NewLogger(cfg.LogLevel)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := db.NewPgxPool(ctx, cfg.DatabaseURL, cfg.MaxDBConns)
	if err != nil {
		return nil, fmt.Errorf("db connect: %w", err)
	}

	userRepo := pgres.NewPostgresUserRepo(pool)
	teamRepo := pgres.NewPostgresTeamRepo(pool, userRepo)
	prRepo := pgres.NewPostgresPRRepo(pool)

	assignSvc := service.NewAssignmentService(userRepo, prRepo)
	teamUC := usecase.NewTeamUsecase(teamRepo, userRepo)
	userUC := usecase.NewUserUsecase(userRepo, teamRepo)
	prUC := usecase.NewPRUsecase(prRepo, userRepo, teamRepo, assignSvc)
	statsUC := usecase.NewStatsUsecase(prRepo)

	mux := http.NewServeMux()
	httpapi.RegisterHandlers(mux, logger, teamUC, userUC, prUC, statsUC, cfg.AdminToken, cfg.UserToken, cfg.AuthorToken, cfg.ReviewerToken)

	mux.HandleFunc("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/app/openapi.yaml")
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_ = httpapi.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	handler := httpapi.LoggingMiddleware(logger)(mux)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	return &App{
		cfg:    cfg,
		logger: logger,
		pool:   pool,
		server: srv,
	}, nil
}


func (a *App) Start(ctx context.Context) error {
	a.logger.Info("starting server", "addr", a.server.Addr)
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Fatal("http server failed", "err", err)
		}
	}()
	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	a.logger.Info("shutting down server")
	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}
	a.pool.Close()
	return nil
}

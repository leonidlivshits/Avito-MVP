package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/leonidlivshits/Avito-MVP/internal/app"
	"github.com/leonidlivshits/Avito-MVP/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	a, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("app init: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := a.Start(ctx); err != nil {
		log.Fatalf("app start: %v", err)
	}

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("app shutdown: %v", err)
	}
}

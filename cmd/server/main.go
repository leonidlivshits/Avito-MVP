package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	port := getEnv("PORT", "8080")
	configPath := getEnv("CONFIG_PATH", "/etc/avito-mvp/config.yaml")
	migrationsDir := getEnv("MIGRATIONS_DIR", "/migrations")

	if data, err := os.ReadFile(configPath); err == nil {
		log.Printf("Loaded config from %s (%d bytes)\n", configPath, len(data))
	} else {
		log.Printf("Config not found at %s: %v\n", configPath, err)
	}

	if files, err := os.ReadDir(migrationsDir); err == nil {
		log.Printf("Migrations dir %s contains %d entries\n", migrationsDir, len(files))
	} else {
		log.Printf("Migrations dir %s not accessible: %v\n", migrationsDir, err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Avito-MVP server. Try /health\n")
	})

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		log.Printf("Starting server on %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server ListenAndServe: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
	log.Println("Server stopped")
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

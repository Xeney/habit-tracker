package main

import (
	"context"
	"habit-tracker/internal/config"
	"habit-tracker/internal/handlers"
	"habit-tracker/internal/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.Load()

	if err := storage.InitDB(cfg.DBPath); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer storage.Close()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", handlers.HealthCheck)

		r.Route("/habits", func(r chi.Router) {
			r.Get("/", handlers.ListHabits)
			r.Post("/", handlers.CreateHabit)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", handlers.GetHabit)
				r.Put("/", handlers.UpdateHabit)
				r.Delete("/", handlers.DeleteHabit)
				r.Post("/log", handlers.LogHabit)
				r.Get("/logs", handlers.GetHabitLogs)
				r.Get("/stats", handlers.GetHabitStats)
			})
		})
	})

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server stopped")
}

package service

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/dupreehkuda/refactor_task/internal/configuration"
	"github.com/dupreehkuda/refactor_task/internal/handlers"
	i "github.com/dupreehkuda/refactor_task/internal/interfaces"
	"github.com/dupreehkuda/refactor_task/internal/storage"
)

type api struct {
	handlers i.Handlers
	config   *configuration.Config
}

func NewByConfig() *api {
	cfg := configuration.New()

	store := storage.New(cfg.FileStoragePath)

	handle := handlers.New(store)

	return &api{
		handlers: handle,
		config:   cfg,
	}
}

// Run runs the service
func (a api) Run() {
	serv := &http.Server{Addr: a.config.Address, Handler: a.service()}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		err := serv.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatalf("Error shutting down: %v", err)
		}
		log.Printf("Server shut down: %s", a.config.Address)
		serverStopCtx()
	}()

	log.Printf("Server started: %s", a.config.Address)
	err := serv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Cant start server: %v", err)
	}

	<-serverCtx.Done()
}

// service returns a custom router
func (a api) service() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(time.Now().String()))
	})

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1/users", func(r chi.Router) {
			r.Get("/", a.handlers.SearchUsers)
			r.Post("/", a.handlers.CreateUser)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", a.handlers.GetUser)
				r.Patch("/", a.handlers.UpdateUser)
				r.Delete("/", a.handlers.DeleteUser)
			})
		})
	})

	return r
}

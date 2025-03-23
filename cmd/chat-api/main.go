package main

import (
	"context"
	"fmt"
	"golang-social-chat/config"
	L "log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Load config file
	LoadConfigS()

	// Connect to Databases
	config.ConnectDB()
	config.ConnectRedis()

	// Define router
	r := chi.NewRouter()

	// Standard middlewares.
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Additional middleware
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.CleanPath)
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.SetHeader("Content-Type", "application/json; charset=utf-8"))

	DefineChatRoutes(r)

	// Create the server
	addr := fmt.Sprintf("%s:%d", Cfg.Address, Cfg.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Channel to listen for OS interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGABRT, syscall.SIGTERM)

	// Running the server in a separate goroutine
	go func() {
		L.Info("Starting server", L.String("address", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			L.Error("Server start error", L.String("error", err.Error()))
		}
	}()

	// Blocking until we receive a signal
	<-stop

	// Creating a deadline for the graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Cfg.ShutdownTimeoutSeconds)*time.Second)
	defer cancel()

	// Attempting a graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		L.Error("Server forced to shutdown",
			L.String("error", err.Error()),
			L.Duration("timeoutSeconds", time.Duration(Cfg.ShutdownTimeoutSeconds)))
	} else {
		L.Info("Server stopped gracefully")
	}
}

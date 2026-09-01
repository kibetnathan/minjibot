package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kibetnathan/minjibot/internal/api"
	"github.com/kibetnathan/minjibot/internal/bot"
)

func main() {
	// Initialize the API first so its /healthz endpoint can come up
	// independently of the Discord bot (which may fail to initialize if
	// DISCORD_TOKEN/DB are not yet reachable).
	apiApp, err := api.NewApp()
	if err != nil {
		log.Fatalf("Failed to initialize API: %v", err)
	}

	// Start the API in a goroutine so the healthcheck responds as soon as
	// possible, regardless of bot state.
	go func() {
		if err := apiApp.Start(); err != nil {
			apiApp.Echo.Logger.Error("API error", "error", err)
		}
	}()

	// Initialize the bot. Failures here (e.g. missing/invalid DISCORD_TOKEN or
	// a database that isn't ready yet) are logged rather than fatal, so the API
	// stays up for healthchecks and can report the service as degraded.
	botApp, err := bot.NewApp()
	if err != nil {
		apiApp.Echo.Logger.Error("Failed to initialize bot (API still running)", "error", err.Error())
	} else {
		go func() {
			if err := botApp.Start(); err != nil {
				botApp.Logger.Error("Bot error", "error", err.Error())
			}
		}()
	}

	// Channel for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	apiApp.Echo.Logger.Info("API is running (bot may be degraded). Press CTRL-C to exit.")

	// Block until interrupt signal
	<-stop
	apiApp.Echo.Logger.Info("Shutting down gracefully...")

	// Shutdown API
	if err := apiApp.Shutdown(context.Background()); err != nil {
		apiApp.Echo.Logger.Error("Error shutting down API", "error", err)
	}

	// Shutdown bot (if it initialized)
	if botApp != nil {
		if err := botApp.Session.Close(); err != nil {
			botApp.Logger.Error("Error closing Discord session", "error", err.Error())
		}
		botApp.Pool.Close()
	}

	apiApp.Pool.Close()
	apiApp.Echo.Logger.Info("Shutdown complete")
}

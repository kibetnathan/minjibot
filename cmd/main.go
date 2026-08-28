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
	// Initialize both applications
	botApp, err := bot.NewApp()
	if err != nil {
		log.Fatalf("Failed to initialize bot: %v", err)
	}

	apiApp, err := api.NewApp()
	if err != nil {
		log.Fatalf("Failed to initialize API: %v", err)
	}

	// Channel for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start bot in goroutine
	go func() {
		if err := botApp.Start(); err != nil {
			botApp.Logger.Error("Bot error", "error", err.Error())
		}
	}()

	// Start API in goroutine
	go func() {
		if err := apiApp.Start(); err != nil {
			apiApp.Echo.Logger.Error("API error", "error", err)
		}
	}()

	botApp.Logger.Info("Both bot and API are running. Press CTRL-C to exit.")

	// Block until interrupt signal
	<-stop
	botApp.Logger.Info("Shutting down gracefully...")

	// Shutdown API
	if err := apiApp.Shutdown(context.Background()); err != nil {
		apiApp.Echo.Logger.Error("Error shutting down API", "error", err)
	}

	// Shutdown bot
	if err := botApp.Session.Close(); err != nil {
		botApp.Logger.Error("Error closing Discord session", "error", err.Error())
	}

	if err := botApp.Conn.Close(context.Background()); err != nil {
		botApp.Logger.Error("Error closing database connection", "error", err)
	}

	botApp.Logger.Info("Shutdown complete")
}
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kibetnathan/minjibot/internal/bot"
)

func main() {
	// Initialize bot application
	botApp, err := bot.NewApp()
	if err != nil {
		log.Fatalf("Failed to initialize bot: %v", err)
	}

	// Channel for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start bot
	go func() {
		if err := botApp.Start(); err != nil {
			botApp.Logger.Error("Bot error", "error", err.Error())
		}
	}()

	botApp.Logger.Info("Bot is running. Press CTRL-C to exit.")

	// Block until interrupt signal
	<-stop
	botApp.Logger.Info("Shutting down gracefully...")

	// Shutdown bot
	if err := botApp.Session.Close(); err != nil {
		botApp.Logger.Error("Error closing Discord session", "error", err.Error())
	}

	botApp.Pool.Close()

	botApp.Logger.Info("Shutdown complete")
}

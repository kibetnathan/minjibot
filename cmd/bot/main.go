package bot

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kibetnathan/minjibot/internal/bot"
)

func main() {
	app, err := bot.NewApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Channel for gracefull shutdown on interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start in goroutine
	go func() {
		if err := app.Start(); err != nil {
			app.Logger.Error("Application error", "error", err.Error())
		}
	}()

	// Block until interrupt signal
	<-stop
	app.Logger.Info("Shutting down bot gracefully...")

	// Perform cleanup manually
	if err := app.Session.Close(); err != nil {
		app.Logger.Error("Error closing Discord session", "error", err.Error())
	}
}

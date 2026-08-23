package api

import (
	"log"
	"os"

	"github.com/kibetnathan/minjibot/internal/api"
)

func main() {
	// Init app
	app, err := api.NewApp()
	if err != nil {
		log.Printf("Failed to initialize application: %v", err)
		os.Exit(1)
	}

	if err := app.Start(); err != nil {
		app.Echo.Logger.Error("Server stopped unexpectedly", "error", err)
		os.Exit(1)
	}
}

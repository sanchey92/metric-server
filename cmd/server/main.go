// Package main is the entry point of the application.
package main

import (
	"context"
	"log"

	"github.com/sanchey92/metric-server/internal/app"
	"github.com/sanchey92/metric-server/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("[main] failed to load config: %v", err)
	}

	ctx := context.Background()

	application, err := app.New(ctx, cfg)
	if err != nil {
		log.Fatalf("[main] failed to initialize app: %v", err)
	}

	if err = application.Run(ctx); err != nil {
		log.Fatalf("[main] application exited with error: %v", err)
	}
}

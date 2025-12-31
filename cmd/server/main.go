package main

import (
	"context"

	"github.com/therenotomorrow/gotes/internal/config"
	"github.com/therenotomorrow/gotes/internal/server"
)

func main() {
	cfg := config.MustNew()
	ctx := context.Background()

	server.Default(cfg).Serve(ctx)
}

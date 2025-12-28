package main

import (
	"github.com/therenotomorrow/gotes/internal/config"
	"github.com/therenotomorrow/gotes/internal/server"
	"github.com/therenotomorrow/gotes/internal/storages/postgres"
	"github.com/therenotomorrow/gotes/pkg/slogx"
)

func main() {
	var (
		cfg = config.MustNew()
		log = slogx.New(slogx.JSON, cfg.Debug)
		app = server.New(cfg, server.Dependencies{
			Database: postgres.MustNew(postgres.Config{
				DSN:    cfg.Postgres.DSN,
				Logger: log,
			}),
			Logger: log,
		})
	)

	defer app.Stop()

	app.Serve()
}

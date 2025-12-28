package config

import (
	"context"
	"time"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
	"github.com/therenotomorrow/ex"
)

const (
	ErrInvalidConfig ex.Error = "invalid config"
)

type Server struct {
	Address  string        `env:"GOTES_SERVER_ADDRESS,required" json:"address"`
	Shutdown time.Duration `                                    json:"shutdown"`
}

type Postgres struct {
	DSN string `env:"GOTES_POSTGRES_DSN,required" json:"dsn"`
}

type Config struct {
	Tier     Tier     `env:"GOTES_TIER,required"  json:"tier"`
	Postgres Postgres `                           json:"postgres"`
	Server   Server   `                           json:"server"`
	Debug    bool     `env:"GOTES_DEBUG,required" json:"debug"`
}

func New(filenames ...string) (Config, error) {
	var cfg Config

	_ = godotenv.Load(filenames...)

	err := envconfig.Process(context.Background(), &cfg)
	if err != nil {
		return Config{}, ErrInvalidConfig.Because(err)
	}

	err = cfg.Tier.Validate()
	if err != nil {
		return Config{}, ErrInvalidConfig.Because(err)
	}

	cfg.Server.Shutdown = 5 * time.Second

	return cfg, nil
}

func MustNew(filenames ...string) Config {
	cfg, err := New(filenames...)

	return ex.Critical(cfg, err)
}

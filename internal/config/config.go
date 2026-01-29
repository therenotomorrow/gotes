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

type MaxConnection struct {
	Idle     time.Duration `json:"idle"`
	Age      time.Duration `json:"age"`
	AgeGrace time.Duration `json:"ageGrace"`
}

type EnforcementPolicy struct {
	MinTime             time.Duration `json:"minTime"`
	PermitWithoutStream bool          `json:"permitWithoutStream"`
}

type KeepAlive struct {
	Time              time.Duration     `json:"time"`
	Timeout           time.Duration     `json:"timeout"`
	MaxConnection     MaxConnection     `json:"maxConnection"`
	EnforcementPolicy EnforcementPolicy `json:"enforcementPolicy"`
}

type Server struct {
	Address              string    `env:"GOTES_SERVER_ADDRESS,required" json:"address"`
	MaxConcurrentStreams uint32    `                                    json:"maxConcurrentStreams"`
	KeepAlive            KeepAlive `                                    json:"keepAlive"`
}

type Postgres struct {
	DSN string `env:"GOTES_POSTGRES_DSN,required" json:"dsn"`
}

type Redis struct {
	Address  string `env:"GOTES_REDIS_ADDRESS,required"  json:"address"`
	Password string `env:"GOTES_REDIS_PASSWORD,required" json:"password"`
}

type Config struct {
	Tier     Tier     `env:"GOTES_TIER,required"  json:"tier"`
	Postgres Postgres `                           json:"postgres"`
	Redis    Redis    `                           json:"redis"`
	Server   Server   `                           json:"server"`
	Debug    bool     `env:"GOTES_DEBUG,required" json:"debug"`
}

func New(filenames ...string) (*Config, error) {
	err := godotenv.Load(filenames...)
	ex.Skip(err)

	cfg := new(Config)

	err = envconfig.Process(context.Background(), cfg)
	if err != nil {
		return nil, ErrInvalidConfig.Because(err)
	}

	err = cfg.Tier.Validate()
	if err != nil {
		return nil, ErrInvalidConfig.Because(err)
	}

	cfg.Server.MaxConcurrentStreams = 100
	cfg.Server.KeepAlive = KeepAlive{
		Time:    30 * time.Second,
		Timeout: 5 * time.Second,
		MaxConnection: MaxConnection{
			Idle:     5 * time.Minute,
			Age:      30 * time.Minute,
			AgeGrace: 5 * time.Second,
		},
		EnforcementPolicy: EnforcementPolicy{
			MinTime:             10 * time.Second,
			PermitWithoutStream: false,
		},
	}

	return cfg, nil
}

func MustNew(filenames ...string) *Config {
	cfg, err := New(filenames...)

	return ex.Critical(cfg, err)
}

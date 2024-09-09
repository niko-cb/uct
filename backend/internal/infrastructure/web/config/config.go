package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	DbUser string `env:"DB_USER" envDefault:"user"`
	DbPass string `env:"DB_PASS,notEmpty"`
	DbPort string `env:"DB_PORT" envDefault:"3306"`
	DbHost string `env:"DB_HOST" envDefault:"localhost"`
	DbName string `env:"DB_NAME" envDefault:"utc"`
	Port   string `env:"PORT" envDefault:"8080"`
}

var Cfg Config // nolint: gochecknoglobals

// Parse parses the environment variables and stores them in the Config struct
func Parse() error {
	if err := env.Parse(&Cfg); err != nil {
		return err
	}
	return nil
}

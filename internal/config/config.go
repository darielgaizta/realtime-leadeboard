package config

import "github.com/caarlos0/env/v11"

type Config struct {
	Port       string `env:"PORT" envDefault:"5000"`
	DBHost     string `env:"DB_HOST" envDefault:"localhost"`
	DBPort     string `env:"DB_PORT" envDefault:"5432"`
	DBName     string `env:"DB_NAME"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	DBUrl      string `env:"DB_URL"`
	JwtSecret  string `env:"JWT_SECRET"`
}

func NewConfiguration() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

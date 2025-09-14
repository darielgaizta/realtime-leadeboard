package config

import "github.com/caarlos0/env/v11"

type Config struct {
	Name             string `env:"NAME" envDefault:"realtime-leaderboard"`
	Port             string `env:"PORT" envDefault:"5000"`
	DBHost           string `env:"DB_HOST" envDefault:"localhost"`
	DBPort           string `env:"DB_PORT" envDefault:"5432"`
	DBName           string `env:"DB_NAME"`
	DBUser           string `env:"DB_USER"`
	DBPassword       string `env:"DB_PASSWORD"`
	DBUrl            string `env:"DB_URL"`
	JWTSecret        string `env:"JWT_SECRET"`
	JWTAccessExpire  int64  `env:"JWT_ACCESS_EXPIRE"`
	JWTRefreshExpire int64  `env:"JWT_REFRESH_EXPIRE"`
	CSRFCookieSecure bool   `env:"CSRF_COOKIE_SECURE"`
}

func NewConfiguration() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

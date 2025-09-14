package middleware

import "github.com/darielgaizta/realtime-leaderboard/internal/config"

type Middleware struct {
	Config *config.Config
}

func NewMiddleware(config *config.Config) *Middleware {
	return &Middleware{
		Config: config,
	}
}

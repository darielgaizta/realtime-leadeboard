package middleware

import "github.com/darielgaizta/realtime-leaderboard/internal/app"

type Middleware struct {
	App *app.App
}

func NewMiddleware(app *app.App) *Middleware {
	return &Middleware{
		App: app,
	}
}

package handler

import (
	"github.com/darielgaizta/realtime-leaderboard/internal/app"
)

type UserHandler struct {
	App *app.App
}

func NewUserHandler(app *app.App) *UserHandler {
	return &UserHandler{App: app}
}

package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/qs-lzh/movie-reservation/internal/app"
)

type MovieHandler struct {
	App *app.App
}

func NewMovieHandler(app *app.App) *MovieHandler {
	return &MovieHandler{
		App: app,
	}
}

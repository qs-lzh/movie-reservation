package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/qs-lzh/movie-reservation/internal/app"
)

type ReservationHandler struct {
	App *app.App
}

func NewReservationHandler(app *app.App) *ReservationHandler {
	return &ReservationHandler{
		App: app,
	}
}

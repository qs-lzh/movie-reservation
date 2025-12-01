package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/qs-lzh/movie-reservation/internal/app"
	"github.com/qs-lzh/movie-reservation/internal/dto"
	"github.com/qs-lzh/movie-reservation/internal/service"
)

type ShowtimeHandler struct {
	App *app.App
}

func NewShowtimeHandler(app *app.App) *ShowtimeHandler {
	return &ShowtimeHandler{
		App: app,
	}
}

// @route GET /api/showtimes/:id
func (h *ShowtimeHandler) GetShowtimeByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		dto.BadRequest(ctx, "Invalid showtime id")
		return
	}
	showtime, err := h.App.ShowtimeService.GetShowtimeByID(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			dto.NotFound(ctx, "Showtime not exists")
			return
		}
		dto.InternalServerError(ctx, "Failed to get showtime")
		return
	}
	dto.Success(ctx, http.StatusOK, showtime)
}

// @route GET /api/showtimes/:id/availability
func (h *ShowtimeHandler) GetShowtimeAvailability(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		dto.BadRequest(ctx, "Invalid showtime id")
		return
	}
	remainingTickets, err := h.App.ReservationService.GetRemainingTickets(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			dto.NotFound(ctx, "Showtime not exists")
			return
		}
		dto.InternalServerError(ctx, "Failed to get showtime availability")
		return
	}
	dto.Success(ctx, http.StatusOK, gin.H{"remaining_tickets": remainingTickets})
}

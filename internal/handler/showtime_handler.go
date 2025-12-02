package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

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

type CreateShowtimeRequest struct {
	MovieID uint      `json:"movie_id" binding:"required"`
	StartAt time.Time `json:"start_at" binding:"required"`
	HallID  uint      `json:"hall_id" binding:"required"`
}

// @route POST /showtimes
func (h *ShowtimeHandler) CreateShowtime(ctx *gin.Context) {
	var req CreateShowtimeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.BadRequest(ctx, "Invalid request body")
		return
	}

	err := h.App.ShowtimeService.CreateShowtime(req.MovieID, req.StartAt, req.HallID)
	if err != nil {
		dto.InternalServerError(ctx, "Failed to create showtime")
		return
	}

	dto.SuccessWithMessage(ctx, http.StatusCreated, nil, "Showtime created successfully")
}

type UpdateShowtimeRequest struct {
	StartAt time.Time `json:"start_at"`
	HallID  uint      `json:"hall_id"`
}

// @route PUT /showtimes/:id
func (h *ShowtimeHandler) UpdateShowtime(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		dto.BadRequest(ctx, "Invalid showtime id")
		return
	}

	var req UpdateShowtimeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.BadRequest(ctx, "Invalid request body")
		return
	}

	// Validate that at least one field is provided for update
	if req.StartAt.IsZero() && req.HallID == 0 {
		dto.BadRequest(ctx, "At least one field (start_at or hall_id) must be provided for update")
		return
	}

	err = h.App.ShowtimeService.UpdateShowtime(uint(id), req.StartAt, req.HallID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			dto.NotFound(ctx, "Showtime not exists")
			return
		}
		dto.InternalServerError(ctx, "Failed to update showtime")
		return
	}

	dto.SuccessWithMessage(ctx, http.StatusOK, nil, "Showtime updated successfully")
}

// @route DELETE /showtimes/:id
func (h *ShowtimeHandler) DeleteShowtimeByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		dto.BadRequest(ctx, "Invalid showtime id")
		return
	}

	err = h.App.ShowtimeService.DeleteShowtimeByID(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			dto.NotFound(ctx, "Showtime not exists")
			return
		}
		dto.InternalServerError(ctx, "Failed to delete showtime")
		return
	}

	dto.SuccessWithMessage(ctx, http.StatusOK, nil, "Showtime deleted successfully")
}

// @route GET /showtimes/:id
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

// @route GET /showtimes/:id/availability
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

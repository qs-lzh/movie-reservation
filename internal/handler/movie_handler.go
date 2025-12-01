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

type MovieHandler struct {
	App *app.App
}

func NewMovieHandler(app *app.App) *MovieHandler {
	return &MovieHandler{
		App: app,
	}
}

// @route GET /api/movies
func (h *MovieHandler) GetAllMovies(ctx *gin.Context) {
	movies, err := h.App.MovieService.GetAllMovies()
	if err != nil {
		dto.InternalServerError(ctx, "Failed to get all movies")
		return
	}
	dto.Success(ctx, http.StatusOK, movies)
}

// @route GET /api/movies/:id
func (h *MovieHandler) GetMovieByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		dto.BadRequest(ctx, "Invalid movie id")
		return
	}
	movie, err := h.App.MovieService.GetMovieByID(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			dto.NotFound(ctx, "Movie not exists")
			return
		}
		dto.InternalServerError(ctx, "failed to get movie")
		return
	}
	dto.Success(ctx, http.StatusOK, movie)
}

func (h *MovieHandler) GetMovieShowtimes(ctx *gin.Context) {
	idParam := ctx.Param("id")
	movieID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		dto.BadRequest(ctx, "Invalid movie id")
		return
	}
	showtimes, err := h.App.ShowtimeService.GetShowtimesByMovieID(uint(movieID))
	if err != nil {
		dto.InternalServerError(ctx, "failed to get showtimes")
		return
	}
	dto.Success(ctx, http.StatusOK, showtimes)
}

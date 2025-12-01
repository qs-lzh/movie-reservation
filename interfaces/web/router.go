package web

import (
	"github.com/gin-gonic/gin"

	"github.com/qs-lzh/movie-reservation/internal/app"
	"github.com/qs-lzh/movie-reservation/internal/handler"
)

func InitRouter(app *app.App) *gin.Engine {
	r := gin.Default()

	authHandler := handler.NewAuthHandler(app)
	movieHandler := handler.NewMovieHandler(app)
	showtimeHandler := handler.NewShowtimeHandler(app)
	reservationHandler := handler.NewReservationHandler(app)

	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	movies := r.Group("movies")
	{
		movies.GET("/", movieHandler.GetAllMovies)
		movies.GET("/:id", movieHandler.GetMovieByID)
		movies.GET("/:id/showtimes", movieHandler.GetMovieShowtimes)
		// [Admin]
		movies.POST("/", movieHandler.CreateMovie)
		movies.PUT("/:id", movieHandler.UpdateMovie)
		movies.DELETE("/:id", movieHandler.DeleteMovie)
	}

	showtimes := r.Group("showtimes")
	{
		showtimes.GET("/:id", showtimeHandler.GetShowtimeByID)
		showtimes.GET("/:id/availability", showtimeHandler.GetShowtimeAvailability)
		// [Admin]
		showtimes.POST("/", showtimeHandler.CreateShowtime)
		showtimes.PUT("/:id", showtimeHandler.UpdateShowtime)
		showtimes.DELETE("/:id", showtimeHandler.DeleteShowtimeByID)
	}

	reservations := r.Group("reservations")
	{
		// [User]
		reservations.POST("/", reservationHandler.CreateReservation)
		reservations.GET("/me", reservationHandler.GetMyReservations)
		reservations.DELETE("/:id", reservationHandler.CancelReservation)
	}

	return r
}

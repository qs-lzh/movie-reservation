package web

import (
	"github.com/gin-gonic/gin"

	"github.com/qs-lzh/movie-reservation/internal/app"
	"github.com/qs-lzh/movie-reservation/internal/handler"
	"github.com/qs-lzh/movie-reservation/internal/middleware"
)

func InitRouter(app *app.App) *gin.Engine {
	authHandler := handler.NewAuthHandler(app)
	movieHandler := handler.NewMovieHandler(app)
	showtimeHandler := handler.NewShowtimeHandler(app)
	reservationHandler := handler.NewReservationHandler(app)

	r := gin.Default()
	// r.Use(gin.Logger())
	// r.Use(gin.Recovery())

	users := r.Group("/users")
	{
		// [User] [Admin]
		users.POST("/register", authHandler.Register)
		users.POST("/login", authHandler.Login)
		users.POST("/logout", authHandler.Logout)
	}

	movies := r.Group("movies")
	movies.Use(middleware.RequireAuth())
	{
		// [User] [Admin]
		movies.GET("/", movieHandler.GetAllMovies)
		movies.GET("/:id", movieHandler.GetMovieByID)
		movies.GET("/:id/showtimes", movieHandler.GetMovieShowtimes)
		// [Admin]
		movies.POST("/", middleware.RequireAdmin(), movieHandler.CreateMovie)
		movies.PUT("/:id", middleware.RequireAdmin(), movieHandler.UpdateMovie)
		movies.DELETE("/:id", middleware.RequireAdmin(), movieHandler.DeleteMovie)
	}

	showtimes := r.Group("showtimes")
	showtimes.Use(middleware.RequireAuth())
	{
		showtimes.GET("/", showtimeHandler.ListAllShowtimes)
		showtimes.GET("/:id", showtimeHandler.GetShowtimeByID)
		showtimes.GET("/:id/availability", showtimeHandler.GetShowtimeAvailability)
		// [Admin]
		showtimes.POST("/", middleware.RequireAdmin(), showtimeHandler.CreateShowtime)
		showtimes.PUT("/:id", middleware.RequireAdmin(), showtimeHandler.UpdateShowtime)
		showtimes.DELETE("/:id", middleware.RequireAdmin(), showtimeHandler.DeleteShowtimeByID)
	}

	reservations := r.Group("reservations")
	reservations.Use(middleware.RequireAuth())
	{
		// [User]
		reservations.POST("/", reservationHandler.CreateReservation)
		reservations.GET("/me", reservationHandler.GetMyReservations)
		reservations.DELETE("/:id", reservationHandler.CancelReservation)
	}

	return r
}

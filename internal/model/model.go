package model

import (
	"time"
)

type User struct {
	ID             uint     `gorm:"primaryKey" json:"id"`
	Name           string   `gorm:"not null;unique" json:"name"`
	HashedPassword string   `gorm:"not null" json:"hashed_password"`
	Role           UserRole `gorm:"not null";not null" json:"role"`
}

type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

type Movie struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Title       string     `gorm:"type:varchar(100);not null" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	Showtimes   []Showtime `gorm:"foreignKey:MovieID" json:"showtimes"`
}

var SeatCount int = 100

type Showtime struct {
	ID      uint      `gorm:"primaryKey" json:"id"`
	MovieID uint      `gorm:"not null" json:"movie_id"`
	StartAt time.Time `gorm:"not null" json:"start_at"`
	HallID  uint      `gorm:"not null" json:"hall_id"`
}

type Reservation struct {
	ID         uint `gorm:"primaryKey" json:"id"`
	ShowtimeID uint `gorm:"not null" json:"showtime_id"`
	UserID     uint `gorm:"not null" json:"user_id"`
}

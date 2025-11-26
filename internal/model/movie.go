package model

import (
	"time"
)

type Movie struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Title       string     `gorm:"type:varchar(100);not null" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	Genre       string     `gorm:"type:varchar(50)" json:"genre"`
	Showtimes   []Showtime `gorm:"foreignKey:MovieID" json:"showtimes"`
}

type Showtime struct {
	ID      uint      `gorm:"primaryKey" json:"id"`
	MovieID uint      `gorm:"not null" json:"movie_id"`
	StartAt time.Time `gorm:"not null" json:"start_at"`
}

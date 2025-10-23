package models

import (
	"time"
)

type Presentation struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	PlayID         uint      `json:"play_id" gorm:"not null;index"`
	Play           Play      `gorm:"foreignKey:PlayID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"play,omitempty"`
	DateTime       time.Time `json:"datetime"`
	Location       string    `json:"location"`
	Price          float64   `json:"price"`
	SeatLimit      int       `json:"seat_limit"`
	AvailableSeats int       `json:"available_seats"`
}

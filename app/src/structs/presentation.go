package structs

import "time"

type Presentation struct {
	ID             uint      `json:"id"`
	PlayID         uint      `json:"play_id"`
	DateTime       time.Time `json:"datetime"`
	Location       string    `json:"location"`
	Price          float64   `json:"price"`
	SeatLimit      int       `json:"seat_limit"`
	AvailableSeats int       `json:"available_seats"`
}

type PresentationsList struct {
	Presentations []Presentation `json:"presentations"`
}

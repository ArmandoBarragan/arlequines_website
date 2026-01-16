package models

import "time"

type Payment struct {
	ID             uint         `gorm:"primaryKey" json:"id"`
	PresentationID uint         `json:"presentation_id" gorm:"not null;index"`
	Presentation   Presentation `gorm:"foreignKey:PresentationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"presentation,omitempty"`
	Amount         float64      `json:"amount"`
	Quantity       int          `json:"quantity"`
	Email          string       `json:"email"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	SessionID      string       `json:"session_id"`
}

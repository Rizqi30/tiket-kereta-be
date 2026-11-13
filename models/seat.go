package models

import "gorm.io/gorm"

type Seat struct {
	gorm.Model
	CarriageID uint `json:"carriage_id" gorm:"not null"`
	SeatNumber string `json:"seat_number" gorm:"not null"` // 1A, 1B, 2A
}
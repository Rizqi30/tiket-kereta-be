package models

import "gorm.io/gorm"

type Carriage struct {
	gorm.Model
	TrainID uint `json:"train_id" gorm:"not null"`
	Name string `json:"name" gorm:"not null"`
	Class string `json:"class" gorm:"not null"` // Eco, Eksekutif, Luxury
	TotalSeats int `json:"total_seats" gorm:"not null"`
	Seats []Seat `json:"seats" gorm:"foreignKey:CarriageID"`
}
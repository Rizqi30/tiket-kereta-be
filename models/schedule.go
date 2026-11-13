package models

import "gorm.io/gorm"

type Schedule struct {
	gorm.Model
	TrainID uint `json:"train_id" gorm:"not null"`
	DepartureStationID uint `json:"departure_station_id" gorm:"not null"`
	ArrivalStationID uint `json:"arrival_station_id" gorm:"not null"`
	DepartureTime string `json:"departure_time" gorm:"not null"`
	ArrivalTime string `json:"arrival_time" gorm:"not null"`
	Date string `json:"date" gorm:"not null"`
	PriceEkonomi int `json:"price_ekonomi"`
	PriceEksekutif int `json:"price_eksekutif"`

	// relasi: GORM otomatis JOIN data kereta & stasiun saat kita pakai Period
	Train Train `json:"train" gorm:"foreignKey:TrainID"`
	DepartureStation Station `json:"departure_station" gorm:"foreignKey:DepartureStationID"`
	ArrivalStation Station `json:"arrival_station" gorm:"foreignKey:ArrivalStationID"`
}
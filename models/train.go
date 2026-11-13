package models

import "gorm.io/gorm"

type Train struct {
	gorm.Model
	Name string `json:"name" gorm:"not null"`
	Code string `json:"code" gorm:"unique;not null"`
	Carriages []Carriage `json:"carriages" gorm:"foreignKey:TrainID"`
}
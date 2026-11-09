package models

import "gorm.io/gorm"

type Station struct {
	gorm.Model
	Name string `json:"name" gorm:"not null"`
	Code string `json:"code" gorm:"unique;not null"`
	City string `json:"city" gorm:"not null"`
}
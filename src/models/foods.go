package models

import "time"

type Food struct {
	Model
	Name           string    `json:"name"`
	Quantity       float64   `json:"quantity"`
	UnitId         uint      `json:"unit_id"`
	FoodUnit       FoodUnit  `gorm:"foreignKey:UnitId"`
	ExpirationDate time.Time `json:"expiration_date"`
	TypeId         uint      `json:"type"`
	FoodType       FoodType  `gorm:"foreignKey:TypeId"`
	UserId         uint      `json:"user_id"`
	User           User      `gorm:"foreignKey:UserId"`
}

type FoodUnit struct {
	Model
	Unit string `json:"unit"`
	Food []Food `gorm:"foreignKey:UnitId"`
}

type FoodType struct {
	Model
	Type string `json:"type"`
	Food []Food `gorm:"foreignKey:TypeId"`
}

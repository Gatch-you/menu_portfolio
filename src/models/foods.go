package models

import "time"

type Food struct {
	Model
	Name           string    `json:"name" gorm:"type:varchar(100);index:idx_member, priority:2"`
	Quantity       float64   `json:"quantity"`
	UnitId         uint      `json:"unit_id"`
	FoodUnit       FoodUnit  `json:"unit_obj" gorm:"foreignKey:UnitId"`
	ExpirationDate time.Time `json:"expiration_date"`
	TypeId         uint      `json:"type_id"`
	FoodType       FoodType  `json:"type" gorm:"foreignKey:TypeId"`
	UserId         uint      `json:"user_id" gorm:"index:idx_member, priority:1"`
	User           User      `json:"-" gorm:"foreignKey:UserId"`
	UseAmount      float64   `json:"use_amount" gorm:"-"`
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

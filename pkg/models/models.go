package model

import "time"

// foods_appにて用いるmodel
type Food struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Quantity       float64   `json:"quantity"`
	Unit           string    `json:"unit"`
	ExpirationDate time.Time `json:"expirationDate"`
	Type           string    `json:"type"`
}

// recipes_appにて用いるmodels

type Recipe struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Image         any    `json:"image"`
	Making_method string `json:"makingMethod"`
}

// recipe_food_appにて用いるmodels
type RecipeFood struct {
	ID                 int     `json:"id"`
	RecipeId           int     `json:"recipe_id"`
	RecipeName         string  `json:"recipe_name"`
	RecipeDescription  string  `json:"recipe_description"`
	FoodId             int     `json:"food_id"`
	FoodName           string  `json:"food_name"`
	UseAmount          float64 `json:"use_amount"`
	FoodUnit           string  `json:"food_unit"`
	RecipeMakingMethod string  `json:"recipe_makingMethod"`
}

type RecipeFoodArray struct {
	FoodID    int     `json:"food_id"`
	RecipeID  int     `json:"recipe_id"`
	UseAmount float64 `json:"use_amount"`
}

type FoodsWithExpiration struct {
	ID             int       `json:"id"`
	FoodId         int       `json:"food_id"`
	FoodName       string    `json:"food_name"`
	FoodQuantity   float64   `json:"food_quantity"`
	FoodUnit       string    `json:"food_unit"`
	ExpirationDate time.Time `json:"expirationDate"`
	RecipeId       int       `json:"recipe_id"`
	RecipeName     string    `json:"recipe_name"`
	UseAmount      float64   `json:"use_amount"`
}

type RecipeID struct {
	RecipeID int `json:"recipe_id"`
}

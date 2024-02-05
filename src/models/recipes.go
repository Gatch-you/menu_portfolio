package models

type Recipe struct {
	Model
	Name         string `json:"name"`
	Description  string `json:"description"`
	MakingMethod string `json:"making_method"`
	UserId       uint   `json:"user_id"`
	Foods        []Food `json:"foods" gorm:"many2many:recipe_food_relations;"`
}

type RecipeFoodRelation struct {
	Model
	FoodId    uint    `json:"food_id"`
	Food      Food    `gorm:"foreignKey:FoodId"`
	RecipeId  uint    `json:"recipe_id"`
	Recipe    Recipe  `gorm:"foreignKey:RecipeId"`
	UseAmount float64 `json:"use_amount"`
}

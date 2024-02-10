package controllers

import (
	database "backend/src/database"
	"backend/src/middlewares"
	"backend/src/models"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func FetchRecipes(c *fiber.Ctx) error {
	var recipes []models.Recipe

	userId, _ := middlewares.GetUserId(c)

	database.DB.Where("user_id = ?", userId).Find(&recipes)

	return c.JSON(recipes)
}

func CreateRecipe(c *fiber.Ctx) error {
	var recipe models.Recipe

	if err := c.BodyParser(&recipe); err != nil {
		return err
	}

	userId, _ := middlewares.GetUserId(c)
	recipe.UserId = userId

	database.DB.Create(&recipe)

	return c.JSON(recipe)
}

func UpdateRecipe(c *fiber.Ctx) error {
	var recipe models.Recipe

	if err := c.BodyParser(&recipe); err != nil {
		return err
	}

	userId, _ := middlewares.GetUserId(c)
	recipe.UserId = userId

	database.DB.Model(&recipe).Updates(&recipe)

	return c.JSON(recipe)

}

func DeleteRecipe(c *fiber.Ctx) error {
	recipeId := c.Params("id")

	var recipe models.Recipe

	database.DB.Where("id = ?", recipeId).Delete(&recipe)

	return c.JSON(fiber.Map{
		"message": "Success",
	})
}

// func FetchRecipeWithFoods(c *fiber.Ctx) error {
// 	recipeId := c.Params("id")

//		// userId, _ := middlewares.GetUserId(c)
//		var recipe models.Recipe
//		fmt.Println(recipeId)
//		// fmt.Println(userId)
//		if err := database.DB.Preload("Foods").Preload("Foods.FoodUnit").Where("id = ? ", recipeId).First(&recipe).Error; err != nil {
//			c.Status(fiber.StatusNotFound)
//			return c.JSON(fiber.Map{
//				"message": "レシピが見つかりません",
//			})
//		}
//		return c.JSON(recipe)
//	}
func FetchRecipeWithFoods(c *fiber.Ctx) error {
	type FoodResponse struct {
		ID             uint      `json:"id"`
		Name           string    `json:"name"`
		Quantity       float64   `json:"quantity"`
		UnitId         uint      `json:"unit_id"`
		Unit           string    `json:"unit"`
		ExpirationDate time.Time `json:"expiration_date"`
		TypeId         uint      `json:"type_id"`
		Type           string    `json:"type"`
		UseAmount      float64   `json:"use_amount"`
		UserId         uint      `json:"user_id"`
	}
	type RecipeResponse struct {
		ID           uint           `json:"id"`
		Name         string         `json:"name"`
		Description  string         `json:"description"`
		MakingMethod string         `json:"making_method"`
		UserId       uint           `json:"user_id"`
		Foods        []FoodResponse `json:"foods"`
	}

	recipeId := c.Params("id")
	fmt.Println(recipeId)

	var recipeResponse RecipeResponse

	if err := database.DB.Model(&models.Recipe{}).Where("id = ?", recipeId).First(&recipeResponse).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "レシピが見つかりません"})
	}

	fmt.Println(recipeResponse)

	var foodResponse []FoodResponse

	if err := database.DB.Table("foods").
		Select("foods.*, recipe_food_relations.use_amount, food_units.unit").
		Joins("join recipe_food_relations on recipe_food_relations.food_id = foods.id").
		Joins("join food_units on food_units.id = foods.unit_id").
		Where("recipe_food_relations.recipe_id = ?", recipeId).
		Scan(&foodResponse).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "関連する食材の取得に失敗しました"})
	}

	recipeResponse.Foods = foodResponse
	return c.JSON(recipeResponse)
}

type FoodToAdd struct {
	FoodId    uint    `json:"food_id"`
	UseAmount float64 `json:"use_amount"`
}

func RegisterFoodToRecipe(c *fiber.Ctx) error {
	recipeId, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.JSON(fiber.Map{
			"message": "URLエンドポイントのパラメータが正しくありません",
		})
	}

	var foodToAdd FoodToAdd
	if err := c.BodyParser(&foodToAdd); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Request body is not valid",
		})
	}

	relation := models.RecipeFoodRelation{
		FoodId:    foodToAdd.FoodId,
		RecipeId:  uint(recipeId),
		UseAmount: foodToAdd.UseAmount,
	}

	if err := database.DB.Create(&relation).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not add food to recipe",
		})
	}

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func DeleteFoodToRecipe(c *fiber.Ctx) error {
	recipeId, err := strconv.ParseUint(c.Params("recipeId"), 10, 64)
	if err != nil {
		return c.JSON(fiber.Map{
			"message": "URLエンドポイント(recipeId)のパラメータが正しくありません",
		})
	}

	foodId, err := strconv.ParseUint(c.Params("foodId"), 10, 64)
	if err != nil {
		return c.JSON(fiber.Map{
			"message": "URLエンドポイント(foodId)のパラメータが正しくありません",
		})
	}

	var deleteFood models.RecipeFoodRelation

	database.DB.Where("recipe_id = ? and food_id = ?", recipeId, foodId).Delete(&deleteFood)

	return c.JSON(fiber.Map{
		"message": "Success",
	})
}

func UpdateFoodToRecipe(c *fiber.Ctx) error {
	recipeId, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.JSON(fiber.Map{
			"message": "URLエンドポイントのパラメータが正しくありません",
		})
	}

	var updateFood models.RecipeFoodRelation
	if err := c.BodyParser(&updateFood); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Request body is not valid",
		})
	}

	database.DB.Where("recipe_id = ? and food_id = ?", recipeId, updateFood.FoodId).Updates(&updateFood)

	return c.JSON(fiber.Map{
		"message": "Success",
	})
}

// // レシピの構造体json形式のデータ変換等も行う

// // レシピの検索に使うが、他の検索アルゴリズムに置き換わる可能性大
// func FetchRecipeByKey(w http.ResponseWriter, r *http.Request) {
// 	db := db.Connect()
// 	defer db.Close()
// }

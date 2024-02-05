package controllers

import (
	database "backend/src/database"
	"backend/src/middlewares"
	"backend/src/models"
	"fmt"
	"strconv"

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
	var recipe models.Recipe

	if err := c.BodyParser(&recipe); err != nil {
		return err
	}

	database.DB.Delete(&recipe)

	return c.JSON(fiber.Map{
		"message": "Success",
	})
}

func FetchRecipeWithFoods(c *fiber.Ctx) error {
	recipeId := c.Params("id")

	// userId, _ := middlewares.GetUserId(c)

	var recipe models.Recipe
	fmt.Println(recipeId)
	// fmt.Println(userId)
	if err := database.DB.Preload("Foods").Preload("Foods.FoodUnit").Where("id = ? ", recipeId).First(&recipe).Error; err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "レシピが見つかりません",
		})
	}

	return c.JSON(recipe)
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
		RecipeId:  uint(recipeId), // recipeID を適切な型に変換する必要があるかもしれません
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

// // レシピの構造体json形式のデータ変換等も行う

// // レシピの検索に使うが、他の検索アルゴリズムに置き換わる可能性大
// func FetchRecipeByKey(w http.ResponseWriter, r *http.Request) {
// 	db := db.Connect()
// 	defer db.Close()
// }

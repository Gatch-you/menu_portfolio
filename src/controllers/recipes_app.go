package controllers

import (
	database "backend/src/database"
	"backend/src/middlewares"
	"backend/src/models"

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

// // レシピの構造体json形式のデータ変換等も行う

// // レシピの検索に使うが、他の検索アルゴリズムに置き換わる可能性大
// func FetchRecipeByKey(w http.ResponseWriter, r *http.Request) {
// 	db := db.Connect()
// 	defer db.Close()
// }

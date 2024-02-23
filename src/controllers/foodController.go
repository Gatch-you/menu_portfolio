package controllers

import (
	database "backend/src/database"
	"backend/src/middlewares"
	"backend/src/models"
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func FetchFoods(c *fiber.Ctx) error {
	var foods []models.Food
	var ctx = context.Background()

	searchWord := c.Query("s")
	userId, _ := middlewares.GetUserId(c)

	if searchWord == "" {
		result, err := database.Cache.Get(ctx, "foods_list_"+strconv.Itoa(int(userId))).Result()

		if err != nil {
			database.DB.Where("user_id = ?", userId).Preload("FoodUnit").Preload("FoodType").Find(&foods)

			bytes, _ := json.Marshal(foods)

			if errKey := database.Cache.Set(ctx, "foods_list_"+strconv.Itoa(int(userId)), bytes, 30*time.Minute).Err(); errKey != nil {
				database.DB.Where("user_id = ?", userId).Preload("FoodUnit").Preload("FoodType").Find(&foods)
			}
		}

		json.Unmarshal([]byte(result), &foods)
	} else {
		database.DB.Where("name LIKE ? AND user_id = ?", "%"+searchWord+"%", userId).Preload("FoodUnit").Preload("FoodType").Find(&foods)
	}

	return c.JSON(foods)
}

func AllFoods(c *fiber.Ctx) error {
	var foods []models.Food
	userId, _ := middlewares.GetUserId(c)

	database.DB.Where("user_id = ?", userId).Preload("FoodUnit").Preload("FoodType").Find(&foods)

	return c.JSON(foods)
}

func CreateFood(c *fiber.Ctx) error {
	var food models.Food

	if err := c.BodyParser(&food); err != nil {
		return err
	}

	userId, _ := middlewares.GetUserId(c)
	food.UserId = userId

	expirationDate, err := time.Parse("2006-01-02", food.ExpirationDate.Format("2006-01-02"))
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "日付の型が間違っています。",
		})
	}
	food.ExpirationDate = expirationDate

	database.DB.Create(&food)

	go database.ClearCache("foods_list_" + strconv.Itoa(int(userId)))

	return c.JSON(food)
}

func UpdateFood(c *fiber.Ctx) error {
	var food models.Food

	if err := c.BodyParser(&food); err != nil {
		return err
	}

	userId, _ := middlewares.GetUserId(c)
	food.UserId = userId

	expirationDate, err := time.Parse("2006-01-02", food.ExpirationDate.Format("2006-01-02"))
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "日付の型が間違っています。",
		})
	}
	food.ExpirationDate = expirationDate

	database.DB.Model(&food).Updates(&food)

	go database.ClearCache("foods_list_" + strconv.Itoa(int(userId)))

	return c.JSON(food)
}

func SoftDeleteFood(c *fiber.Ctx) error {
	foodId := c.Params("id")
	userId, _ := middlewares.GetUserId(c)

	var food models.Food

	if err := c.BodyParser(&food); err != nil {
		return err
	}

	if err := database.DB.
		Where("id = ? AND user_id = ?", foodId, userId).
		First(&food).Error; err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "食材が見つかりません。",
		})
	}

	database.DB.Model(&food).Updates(map[string]interface{}{
		"quantity": 0.0,
	})

	go database.ClearCache("foods_list_" + strconv.Itoa(int(userId)))

	return c.JSON(food)
}

func DeleteFood(c *fiber.Ctx) error {
	var food models.Food
	userId, _ := middlewares.GetUserId(c)

	if err := c.BodyParser(&food); err != nil {
		return err
	}

	database.DB.Delete(&food)

	go database.ClearCache("foods_list_" + strconv.Itoa(int(userId)))

	return c.JSON(fiber.Map{
		"message": "Success",
	})
}

func FetchFoodswithExpiration(c *fiber.Ctx) error {
	expirationDate := c.Query("expiration_date")
	userId, _ := middlewares.GetUserId(c)

	var foodResponse []models.Food

	// まずは期限の切れる食材を取得
	if err := database.DB.Model(&models.Food{}).
		Where("user_id = ? AND expiration_date <= DATE_ADD(DATE(NOW()), INTERVAL ? DAY)", userId, expirationDate).
		Preload("FoodUnit").
		Find(&foodResponse).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "食材が見つかりません"})
	}

	// 事前に取得したfoodに対してrecipe_food_relationを参照し、関連するメニューを取得
	for i, food := range foodResponse {
		var recipes []models.Recipe
		if err := database.DB.Table("recipes").
			Select("recipes.*, recipe_food_relations.use_amount").
			Joins("JOIN recipe_food_relations on recipe_food_relations.recipe_id = recipes.id").
			Where("recipe_food_relations.food_id = ?", food.Id).
			Scan(&recipes).Error; err != nil {
			continue
		}
		foodResponse[i].Recipes = recipes
	}
	return c.JSON(foodResponse)
}

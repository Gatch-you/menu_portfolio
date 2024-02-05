package controllers

import (
	database "backend/src/database"
	"backend/src/middlewares"
	"backend/src/models"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func FetchFoods(c *fiber.Ctx) error {
	var foods []models.Food

	user_id, _ := middlewares.GetUserId(c)

	database.DB.Where("user_id = ?", user_id).Find(&foods)

	return c.JSON(foods)
}

// todo:バリテーションの実施
func CreateFoods(c *fiber.Ctx) error {
	var food models.Food

	if err := c.BodyParser(&food); err != nil {
		return err
	}

	user_id, _ := middlewares.GetUserId(c)
	food.UserId = user_id

	expirationDate, err := time.Parse("2006-01-02", food.ExpirationDate.Format("2006-01-02"))
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "日付の型が間違っています。",
		})
	}
	food.ExpirationDate = expirationDate

	database.DB.Create(&food)

	return c.JSON(food)
}

func UpdateFood(c *fiber.Ctx) error {
	var food models.Food

	if err := c.BodyParser(&food); err != nil {
		return err
	}

	user_id, _ := middlewares.GetUserId(c)
	food.UserId = user_id

	expirationDate, err := time.Parse("2006-01-02", food.ExpirationDate.Format("2006-01-02"))
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "日付の型が間違っています。",
		})
	}
	food.ExpirationDate = expirationDate

	database.DB.Model(&food).Updates(&food)

	return c.JSON(food)
}

func SoftDeleteFoods(c *fiber.Ctx) error {
	var food models.Food

	if err := c.BodyParser(&food); err != nil {
		return err
	}

	userId, _ := middlewares.GetUserId(c)
	foodId := food.Id

	if err := database.DB.Where("id = ? AND user_id = ?", foodId, userId).First(&food).Error; err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "食材が見つかりません。",
		})
	}

	fmt.Println(food)

	database.DB.Model(&food).Updates(map[string]interface{}{
		"quantity": 0.0,
	})

	return c.JSON(food)
}

func DeleteFoods(c *fiber.Ctx) error {
	var food models.Food

	if err := c.BodyParser(&food); err != nil {
		return err
	}

	database.DB.Delete(&food)

	return c.JSON(fiber.Map{
		"message": "Success",
	})
}

// // curl http://localhost:8080/backend/foods/search_name%q=
// func SearchFoodsName(w http.ResponseWriter, r *http.Request) {
// 	query := r.URL.Query().Get("q")

// 	db := db.Connect()
// 	defer db.Close()

// 	rows, err := db.Query("SELECT name, quantity, unit, expiration_date FROM foods WHERE name LIKE ?", "%"+query+"%")
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}

// 	foodArgs := make([]model.Food, 0)
// 	for rows.Next() {
// 		var food model.Food
// 		err = rows.Scan(&food.Name, &food.Quantity, &food.Unit, &food.ExpirationDate)
// 		if err != nil {
// 			log.Fatal(err.Error())
// 		}
// 		foodArgs = append(foodArgs, food)
// 	}

// 	middleware.CorsMiddleware(http.DefaultServeMux)
// 	json.NewEncoder(w).Encode(foodArgs)

// }

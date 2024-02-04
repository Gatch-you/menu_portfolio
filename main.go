package main

import (
	database "backend/src/database"
	"backend/src/routes"

	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	database.Connect()
	database.AutoMigrate()
	log.Println("Hello, I'm Menu Proposer!")

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://localhost:3000",
	}))

	routes.Setup(app)

	app.Listen(":3000")

	// //controllersのリクエスト処理
	// http.HandleFunc("/backend/test", connection_test)
	// http.HandleFunc("/backend/foods", controllers.FetchFoods)
	// http.HandleFunc("/backend/insert_food", controllers.InsertFoods)
	// http.HandleFunc("/backend/delete_food", controllers.DeleteFoods)
	// http.HandleFunc("/backend/update_food", controllers.UpdateFoods)
	// http.HandleFunc("/backend/search_name", controllers.SearchFoodsName)

	// //controllersのリクエスト
	// http.HandleFunc("/backend/recipes", controllers.FetchRecipes)
	// http.HandleFunc("/backend/insert_recipe", controllers.InsertRecipe)
	// http.HandleFunc("/backend/update_recipe", controllers.UpdateRecipe)
	// http.HandleFunc("/backend/delete_recipe", controllers.DeleteRecipe)

	// //recipe_foodのリクエスト
	// http.HandleFunc("/backend/recipe_food", controllers.FetchRecipesWithFood)
	// http.HandleFunc("/backend/recipe_food/update_food_storage/", controllers.UpdateFoodStorage)
	// http.HandleFunc("/backend/recipe_food/update_using_food_quantity", controllers.UpdateUsingFoodQuantity)
	// http.HandleFunc("/backend/recipe_food/insert_use_food_array", controllers.InsertUseFoodArray)
	// http.HandleFunc("/backend/recipes/", controllers.FetchRecipeDetail)
	// http.HandleFunc("/backend/delete_using_food", controllers.DeleteUsingFood)
	// http.HandleFunc("/backend/recipe_food/insert_use_food", controllers.InsertUseFood)
	// http.HandleFunc("/backend/recipe_food/foods_expiration", controllers.ShowFoodsWithExpiration)
	// //goroutineで定時に発火させる→今後Line等による機能の実装を可能にする。
	// go controllers.FetchExpirationFood(nil, nil)

	// // CORSミドルウェア関数を用いてCORS解決
	// CORSHandler := middleware.CorsMiddleware(http.DefaultServeMux)
	// http.ListenAndServe(":8000", CORSHandler)
}

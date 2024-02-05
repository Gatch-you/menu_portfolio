package routes

import (
	"backend/src/controllers"
	"backend/src/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	api := app.Group("api")
	api.Get("test", connection_test)

	user := api.Group("user")
	tmp := api.Group("tmp")
	user.Post("register", controllers.Resister)
	user.Post("login", controllers.Login)
	tmp.Post("password/reset/request", controllers.PasswordResetRequest)

	// // パスワードを忘れた人へのリセット処理
	tmpAuthenticated := tmp.Use(middlewares.IsTmpAuthenticated)
	tmpAuthenticated.Post("password/reset/confirm", controllers.ResetPassword)

	userAuthenticated := user.Use(middlewares.IsAuthenticated)
	userAuthenticated.Post("logout", controllers.Logout)
	userAuthenticated.Get("foods", controllers.FetchFoods)
	userAuthenticated.Post("foods", controllers.CreateFood)
	userAuthenticated.Put("foods", controllers.UpdateFood)
	userAuthenticated.Put("foods/sfdelete", controllers.SoftDeleteFood)
	userAuthenticated.Delete("foods", controllers.DeleteFood)

	userAuthenticated.Get("recipes", controllers.FetchRecipes)
	userAuthenticated.Post("recipes", controllers.CreateRecipe)
	userAuthenticated.Put("recipes", controllers.UpdateRecipe)
	userAuthenticated.Delete("recipes", controllers.DeleteRecipe)
	userAuthenticated.Delete("recipes", controllers.DeleteRecipe)

	userAuthenticated.Get("recipes/detail/:id", controllers.FetchRecipeWithFoods)
	userAuthenticated.Post("recipes/resistfood/:id", controllers.RegisterFoodToRecipe)

	// userAuthenticated := api.Use(middlewares.IsUser)

	// admin := app.Group("admin")

}

func connection_test(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Hello, Golang with fiber",
	})
}

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
	userAuthenticated.Get("profile", controllers.User)
	userAuthenticated.Put("profile", controllers.UpdateInfo)
	userAuthenticated.Put("profile/password", controllers.UpdatePassword)
	userAuthenticated.Post("logout", controllers.Logout)

	userAuthenticated.Get("foods", controllers.FetchFoods)
	userAuthenticated.Post("foods", controllers.CreateFood)
	userAuthenticated.Put("foods", controllers.UpdateFood)
	userAuthenticated.Put("foods/sfdelete/:id", controllers.SoftDeleteFood)
	userAuthenticated.Delete("foods", controllers.DeleteFood)
	userAuthenticated.Get("foods/select", controllers.AllFoods)

	userAuthenticated.Get("recipes", controllers.FetchRecipes)
	userAuthenticated.Post("recipes", controllers.CreateRecipe)
	userAuthenticated.Put("recipes", controllers.UpdateRecipe)
	userAuthenticated.Delete("recipes/:id", controllers.DeleteRecipe)

	userAuthenticated.Get("recipes/detail/:id", controllers.FetchRecipeWithFoods)
	userAuthenticated.Post("recipes/detail/controllfood/:id", controllers.RegisterFoodToRecipe)
	userAuthenticated.Delete("recipes/detail/controllfood/:recipeId/:foodId", controllers.DeleteFoodToRecipe)
	userAuthenticated.Put("recipes/detail/controllfood/:id", controllers.UpdateFoodToRecipe)
	userAuthenticated.Put("recipes/cooking/:id", controllers.MakeDish)

	userAuthenticated.Get("foods/expiration", controllers.FetchFoodswithExpiration)

	// userAuthenticated := api.Use(middlewares.IsUser)
	// admin := app.Group("admin")

}

func connection_test(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Hello, Golang with fiber",
	})
}

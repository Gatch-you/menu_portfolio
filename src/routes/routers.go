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
	userAuthenticated.Post("foods", controllers.CreateFoods)
	userAuthenticated.Put("foods", controllers.UpdateFood)
	userAuthenticated.Put("foods/sfdelete", controllers.SoftDeleteFoods)
	userAuthenticated.Delete("foods", controllers.DeleteFoods)

	// userAuthenticated := api.Use(middlewares.IsUser)

	// admin := app.Group("admin")

}

func connection_test(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Hello, Golang with fiber",
	})
}

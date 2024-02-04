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
	user.Post("register", controllers.Resister)
	user.Post("login", controllers.Login)

	userAuthenticated := user.Use(middlewares.IsAuthenticated)
	userAuthenticated.Get("foods", controllers.Foods)

	// userAuthenticated := api.Use(middlewares.IsUser)

	// admin := app.Group("admin")

}

func connection_test(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Hello, Golang with fiber",
	})
}

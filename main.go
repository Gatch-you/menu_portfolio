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
	database.SetRedis()
	database.SetupCacheChannel()

	log.Println("Hello, I'm Menu Proposer!")

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://localhost:3000",
	}))

	routes.Setup(app)

	app.Listen(":3000")

}

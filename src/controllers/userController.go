package controllers

import (
	database "backend/src/database"
	middlewares "backend/src/middlewares"
	models "backend/src/models"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Resister(c *fiber.Ctx) error {
	var regist_data map[string]string

	if err := c.BodyParser(&regist_data); err != nil {
		fmt.Println(err)
		return err
	}

	if regist_data["password"] != regist_data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "パスワードが間違っています \n 正しいパスワードを入力してください。",
		})
	}

	user := models.User{
		FirstName: regist_data["first_name"],
		LastName:  regist_data["last_name"],
		Email:     regist_data["email"],
		IsAdmin:   strings.Contains(c.Path(), "/api/admin"),
	}

	user.SetPassword(regist_data["password"])

	database.DB.Create(&user)

	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	var login_data map[string]string

	if err := c.BodyParser(&login_data); err != nil {
		return err
	}

	var user models.User

	database.DB.Where("email = ?", login_data["email"]).First(&user)

	if user.Id == 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "ユーザーが見つかりません。\nメールアドレスを確認の上、ログインしてください。",
		})
	}

	if err := user.ComparePassword(login_data["password"]); err != nil {
		return c.JSON(fiber.Map{
			"message": "パスワードが間違っています。",
		})
	}

	isUser := strings.Contains(c.Path(), "/api/user")

	var scope string

	if isUser {
		scope = "user"
	} else {
		scope = "admin"
	}

	if !isUser && user.IsAdmin {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "こちらはadminアカウントです。\n専用のフォームからアクセスしてください。",
		})
	}

	token, err := middlewares.GenerateJWT(user.Id, scope)

	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "許可されたアカウントではありません。",
		})
	}

	cookie := fiber.Cookie{
		Name:     "SID_MCB",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 48),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

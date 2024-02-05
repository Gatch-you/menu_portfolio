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

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "SID_MCB",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "Success",
	})
}

func UpdateInfo(c *fiber.Ctx) error {
	var user_data map[string]string
	if err := c.BodyParser(&user_data); err != nil {
		return err
	}

	id, _ := middlewares.GetUserId(c)

	user := models.User{
		FirstName: user_data["first_name"],
		LastName:  user_data["last_name"],
		Email:     user_data["email"],
	}
	user.Id = id

	database.DB.Model(&user).Updates(&user)

	return c.JSON(user)
}

// ログインをしている状態におけるパスワードの変更機能
func UpdatePassword(c *fiber.Ctx) error {
	var update_data map[string]string

	if err := c.BodyParser(&update_data); err != nil {
		return err
	}

	if update_data["password"] != update_data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "passwords do not match",
		})
	}

	id, _ := middlewares.GetUserId(c)

	user := models.User{}
	user.Id = id

	user.SetPassword(update_data["password"])

	database.DB.Model(&user).Updates(&user)

	fmt.Printf("Changing password is successed!")
	return c.JSON(user)
}

func PasswordResetRequest(c *fiber.Ctx) error {
	var user_request_data map[string]string

	if err := c.BodyParser(&user_request_data); err != nil {
		return err
	}

	email := user_request_data["email"]

	var user models.User

	database.DB.Where("email = ?", email).First(&user)

	if user.Id == 0 {
		return c.JSON(fiber.Map{
			"message": "指定のユーザーアカウントが見つかりません",
		})
	}

	tmpToken, _ := middlewares.GenerateJWT(user.Id, "user")

	cookie := fiber.Cookie{
		Name:     "SID_MCB_TMP",
		Value:    tmpToken,
		Expires:  time.Now().Add(time.Minute * 30),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func ResetPassword(c *fiber.Ctx) error {
	var update_data map[string]string

	if err := c.BodyParser(&update_data); err != nil {
		return err
	}

	if update_data["password"] != update_data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "passwords do not match",
		})
	}

	id, _ := middlewares.GetUserId(c)

	user := models.User{}
	user.Id = id

	user.SetPassword(update_data["password"])

	database.DB.Model(&user).Updates(&user)

	fmt.Printf("Changing password is successed!")
	return c.JSON(user)
}

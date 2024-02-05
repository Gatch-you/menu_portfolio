package middlewares

import (
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

const SecretKey = "secret"

type ClaimsWithScope struct {
	jwt.StandardClaims
	Scope string
}

func IsUser(c *fiber.Ctx) error {
	cookie := c.Cookies("SID_MCB")

	token, err := jwt.ParseWithClaims(cookie, &ClaimsWithScope{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil || !token.Valid {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthorized user",
		})
	}

	payload := token.Claims.(*ClaimsWithScope)
	isAmbassador := strings.Contains(c.Path(), "/api/admin")

	if (payload.Scope == "admin" && isAmbassador) || (payload.Scope == "ambassador" && !isAmbassador) {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthorized",
		})
	}

	return c.Next()
}

func GenerateJWT(id uint, scope string) (string, error) {
	payload := ClaimsWithScope{}

	payload.Subject = strconv.Itoa(int(id))
	payload.ExpiresAt = time.Now().Add(time.Hour * 24).Unix()
	payload.Scope = scope

	return jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte(SecretKey))
}

func IsAuthenticated(c *fiber.Ctx) error {
	cookie := c.Cookies("SID_MCB")

	// cookieからvalueを取得して配列を返す
	token, err := jwt.ParseWithClaims(cookie, &ClaimsWithScope{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	// 取得した配列に対して
	if err != nil || !token.Valid {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "許可されたユーザーでありません、もしくはログイン期限が切れています。",
		})
	}

	payload := token.Claims.(*ClaimsWithScope)
	isUser := strings.Contains(c.Path(), "/api/user")

	if (payload.Scope == "admin" && isUser) || (payload.Scope == "ambassador" && !isUser) {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "許可されたユーザーではありません",
		})
	}

	return c.Next()
}

func IsTmpAuthenticated(c *fiber.Ctx) error {
	cookie := c.Cookies("SID_MCB_TMP")

	// cookieからvalueを取得して配列を返す
	token, err := jwt.ParseWithClaims(cookie, &ClaimsWithScope{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	// 取得した配列に対して
	if err != nil || !token.Valid {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Cookieの有効期間が切れています。再度メールアドレス認証からやり直してください。",
		})
	}

	payload := token.Claims.(*ClaimsWithScope)
	isUser := strings.Contains(c.Path(), "/api/user")

	if (payload.Scope == "admin" && isUser) || (payload.Scope == "ambassador" && !isUser) {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "許可されたユーザーではありません",
		})
	}

	return c.Next()
}

func GetUserId(c *fiber.Ctx) (uint, error) {
	cookie := c.Cookies("SID_MCB")

	token, err := jwt.ParseWithClaims(cookie, &ClaimsWithScope{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil {
		return 0, err
	}

	payload := token.Claims.(*ClaimsWithScope)

	id, _ := strconv.Atoi(payload.Subject)

	return uint(id), nil
}

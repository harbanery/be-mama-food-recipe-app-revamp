package middleware

import (
	"fmt"
	"mama-recipe/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func NewAuthMiddleware(env *helper.EnvLoad) fiber.Handler {
	secretKey := env.JWT_SECRET_KEY
	return func(c *fiber.Ctx) error {
		response := &helper.WebResponse[interface{}]{}

		tokenAuth := c.Get("Authorization", "NOT FOUND")[7:]
		if tokenAuth == "" {
			response = helper.Response(c, 404, "Unathorized", nil)

			c.JSON(response)
		}

		token, err := jwt.Parse(tokenAuth, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})
		if err != nil {
			response = helper.Response(c, 404, "Unathorized", nil)

			c.JSON(response)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response = helper.Response(c, 404, "Unathorized", nil)

			c.JSON(response)
		}

		c.Locals("id", claims["id"].(string))
		c.Locals("email", claims["email"].(string))

		return c.Next()

	}
}

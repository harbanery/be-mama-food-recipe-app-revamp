package middleware

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/microcosm-cc/bluemonday"
)

func NewBodyMiddleware(policy *bluemonday.Policy) fiber.Handler {
	return func(c *fiber.Ctx) error {
		query := c.Request().URI().QueryArgs()
		query.VisitAll(func(key, value []byte) {
			c.Context().QueryArgs().SetBytesKV(key, []byte(policy.Sanitize(string(value))))
		})

		if !c.Is("application/json") {
			return c.Next()
		}

		var sanitized map[string]interface{}
		if err := c.BodyParser(&sanitized); err == nil {
			for k, v := range sanitized {
				if str, ok := v.(string); ok {
					sanitized[k] = policy.Sanitize(str)
				}
			}

			jsonData, err := json.Marshal(sanitized)
			if err != nil {
				return err
			}
			c.Request().SetBody(jsonData)
		}

		return c.Next()
	}
}

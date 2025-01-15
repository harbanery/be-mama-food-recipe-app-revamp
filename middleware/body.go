package middleware

import (
	"encoding/json"
	"mama-recipe/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/microcosm-cc/bluemonday"
)

func NewBodyMiddleware(policy *bluemonday.Policy) fiber.Handler {
	return func(c *fiber.Ctx) error {
		query := c.Request().URI().QueryArgs()
		query.VisitAll(func(key, value []byte) {
			c.Context().QueryArgs().SetBytesKV(key, []byte(policy.Sanitize(string(value))))
		})

		if c.Is("application/json") {
			var sanitized map[string]interface{}
			if err := c.BodyParser(&sanitized); err == nil {
				for k, v := range sanitized {
					if str, ok := v.(string); ok {
						sanitized[k] = policy.Sanitize(str)
					}
				}

				jsonData, err := json.Marshal(sanitized)
				if err != nil {
					response := helper.Response(c, 400, err.Error(), nil)

					return c.JSON(response)
				}
				c.Request().SetBody(jsonData)
			}
		} else if c.Is("multipart/form-data") {
			form, err := c.MultipartForm()
			if err == nil && form != nil {
				sanitized := make(map[string]interface{})

				for key, values := range form.Value {
					for i, value := range values {
						form.Value[key][i] = policy.Sanitize(value)
					}
					sanitized[key] = form.Value[key]
				}
			}
		}

		return c.Next()
	}
}

package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
)

func NewFiber(env *EnvLoad) *fiber.App {
	var app = fiber.New(fiber.Config{
		AppName:      env.APP_NAME,
		ErrorHandler: NewErrorHandler(),
		Prefork:      env.WEB_PREFORK == "true",
	})

	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     env.CORS_ALLOW_ORIGINS,
		AllowMethods:     env.CORS_ALLOW_METHODS,
		AllowHeaders:     env.CORS_ALLOW_HEADERS,
		ExposeHeaders:    env.CORS_EXPOSE_HEADERS,
	}))

	return app
}

func NewErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		return ctx.Status(code).JSON(fiber.Map{
			"errors": err.Error(),
		})
	}
}
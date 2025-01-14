package config

import (
	"mama-recipe/app/account"
	"mama-recipe/app/authentication"
	"mama-recipe/helper"
	"mama-recipe/middleware"
	"mama-recipe/route"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/microcosm-cc/bluemonday"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB          *gorm.DB
	App         *fiber.App
	Validate    *validator.Validate
	Policy      *bluemonday.Policy
	Environment *helper.EnvLoad
}

type Information struct {
	URL           string  `json:"url"`
	Path          string  `json:"path"`
	VersionNumber *string `json:"version_number,omitempty"`
	Info          *string `json:"info,omitempty"`
	Message       string  `json:"message"`
}

func Bootstrap(config *BootstrapConfig) {
	message := "Server is running"

	config.App.Get("/v1", func(c *fiber.Ctx) error {
		version := "v1"

		res := &Information{
			URL:           c.Request().URI().String(),
			Path:          c.Path(),
			VersionNumber: &version,
			Message:       message,
		}

		return c.JSON(res)
	})

	authRepository := authentication.NewAuthRepository()
	authUseCase := authentication.NewAuthUseCase(config.DB, config.Validate, authRepository, config.Environment)
	authHandler := authentication.NewAuthHandler(authUseCase)

	accountRepository := account.NewAccountRepository()
	accountUseCase := account.NewAccountUseCase(config.DB, config.Validate, accountRepository, config.Environment)
	accountHandler := account.NewAccountHandler(accountUseCase)

	authMiddleware := middleware.NewAuthMiddleware(config.Environment)
	bodyMiddleware := middleware.NewBodyMiddleware(config.Policy)

	route := route.RouteConfig{
		App:            config.App,
		AuthMiddleware: authMiddleware,
		BodyMiddleware: bodyMiddleware,
		AuthHandler:    authHandler,
		AccountHandler: accountHandler,
	}

	route.Setup()

	config.App.Get("/", func(c *fiber.Ctx) error {
		info := "Hello, Welcome to API Mama Recipe"

		res := &Information{
			URL:     c.Request().URI().String(),
			Path:    c.Path(),
			Info:    &info,
			Message: message,
		}

		return c.JSON(res)
	})

	config.App.Get("*", func(c *fiber.Ctx) error {
		res := &Information{
			URL:     c.Request().URI().String(),
			Path:    c.Path(),
			Message: "Route not found",
		}

		return c.Status(fiber.StatusNotFound).JSON(res)
	})
}

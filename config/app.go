package config

import (
	"mama-recipe/app/authentication/signin"
	"mama-recipe/app/authentication/signout"
	"mama-recipe/app/authentication/signup"
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

	signUpRepository := signup.NewSignUpRepository()
	signUpUseCase := signup.NewSignUpUseCase(config.DB, config.Validate, signUpRepository, config.Environment)
	signUpHandler := signup.NewSignUpHandler(signUpUseCase)

	signInRepository := signin.NewSignInRepository()
	signInUseCase := signin.NewSignInUseCase(config.DB, config.Validate, signInRepository, config.Environment)
	signInHandler := signin.NewSignInHandler(signInUseCase)

	signOutRepository := signout.NewSignOutRepository()
	signOutUseCase := signout.NewSignOutUseCase(config.DB, config.Validate, signOutRepository, config.Environment)
	signOutHandler := signout.NewSignOutHandler(signOutUseCase)

	authMiddleware := middleware.NewAuthMiddleware(config.Environment)
	bodyMiddleware := middleware.NewBodyMiddleware(config.Policy)

	route := route.RouteConfig{
		App:            config.App,
		AuthMiddleware: authMiddleware,
		BodyMiddleware: bodyMiddleware,
		SignUpHandler:  signUpHandler,
		SignInHandler:  signInHandler,
		SignOutHandler: signOutHandler,
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

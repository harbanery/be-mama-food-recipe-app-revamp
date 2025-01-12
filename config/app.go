package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Validate *validator.Validate
}

type Information struct {
	URL       string `json:"url"`
	Path      string `json:"path"`
	VersionNumber *string `json:"version_number,omitempty"`
	Info      *string `json:"info,omitempty"`
	Message   string `json:"message"`
}

func Bootstrap(config *BootstrapConfig) {
	config.App.Get("/v1", func(c *fiber.Ctx) error {
		version := "v1"

		res := &Information{
			URL: c.Request().URI().String(),
			Path: c.Path(),
			VersionNumber: &version,
			Message: "Server is running",
		}

		return c.JSON(res)
	})

	config.App.Get("*", func(c *fiber.Ctx) error {
		info := "Hello, Welcome to API Mama Recipe"

		res := &Information{
			URL: c.Request().URI().String(),
			Path: c.Path(),
			Info: &info,
			Message: "Server is running",
		}

		return c.JSON(res)
	})
}
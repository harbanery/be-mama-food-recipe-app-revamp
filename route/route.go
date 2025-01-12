package route

import "github.com/gofiber/fiber/v2"

type RouteConfig struct {
	App *fiber.App
}

func (c *RouteConfig) Setup()  {
	c.SignUp()
}

func (c *RouteConfig) SignUp() {
	c.App.Post("/signup")
}
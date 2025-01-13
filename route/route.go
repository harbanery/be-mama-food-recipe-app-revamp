package route

import (
	"mama-recipe/app/authentication/signin"
	"mama-recipe/app/authentication/signout"
	"mama-recipe/app/authentication/signup"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App            *fiber.App
	AuthMiddleware fiber.Handler
	BodyMiddleware fiber.Handler
	SignUpHandler  signup.SignUpHandler
	SignInHandler  signin.SignInHandler
	SignOutHandler signout.SignOutHandler
}

func (c *RouteConfig) Setup() {
	c.SignUpRoute()
	c.SignInRoute()
	c.SignOutRoute()
}

func (c *RouteConfig) SignUpRoute() {
	c.App.Post("/signup", c.BodyMiddleware, c.SignUpHandler.Register)
}

func (c *RouteConfig) SignInRoute() {
	c.App.Post("/signin", c.BodyMiddleware, c.SignInHandler.Login)
}

func (c *RouteConfig) SignOutRoute() {
	c.App.Get("/signout", c.AuthMiddleware, c.SignOutHandler.Logout)
}

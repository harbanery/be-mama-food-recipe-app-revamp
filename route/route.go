package route

import (
	"mama-recipe/app/account"
	"mama-recipe/app/authentication"
	"mama-recipe/app/recipe"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App            *fiber.App
	AuthMiddleware fiber.Handler
	BodyMiddleware fiber.Handler
	AuthHandler    authentication.AuthHandler
	AccountHandler account.AccountHandler
	RecipeHandler  recipe.RecipeHandler
}

func (c *RouteConfig) Setup() {
	c.AuthRoute()
	c.AccountRoute()
	c.RecipeRoute()
}

func (c *RouteConfig) AuthRoute() {
	c.App.Post("/signup", c.BodyMiddleware, c.AuthHandler.Register)
	c.App.Post("/signin", c.BodyMiddleware, c.AuthHandler.Login)
	c.App.Get("/signout", c.AuthMiddleware, c.AuthHandler.Logout)
}

func (c *RouteConfig) AccountRoute() {
	account := c.App.Group("/account")
	account.Get("/profile", c.AuthMiddleware, c.AccountHandler.DetailProfile)
	account.Put("/profile", c.AuthMiddleware, c.BodyMiddleware, c.AccountHandler.UpdateProfile)
	account.Put("/profile-photo", c.AuthMiddleware, c.BodyMiddleware, c.AccountHandler.UpdateProfilePhoto)
}

func (c *RouteConfig) RecipeRoute() {
	recipe := c.App.Group("/recipe")
	recipe.Post("/add", c.AuthMiddleware, c.BodyMiddleware, c.RecipeHandler.AddRecipe)
	recipe.Get("/list", c.AuthMiddleware, c.RecipeHandler.ListRecipe)
}

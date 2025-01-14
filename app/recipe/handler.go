package recipe

import "github.com/gofiber/fiber/v2"

type RecipeHandler interface {
	AddRecipe(ctx *fiber.Ctx) error
	ListRecipe(ctx *fiber.Ctx) error
}

type recipeHandler struct {
	UseCase RecipeUseCase
}

func NewRecipeHandler(useCase RecipeUseCase) RecipeHandler {
	return &recipeHandler{
		UseCase: useCase,
	}
}

func (c *recipeHandler) AddRecipe(ctx *fiber.Ctx) error {
	res := c.UseCase.AddRecipe(ctx)
	return ctx.Status(res.Code).JSON(res)
}

func (c *recipeHandler) ListRecipe(ctx *fiber.Ctx) error {
	res := c.UseCase.ListRecipe(ctx)
	return ctx.Status(res.Code).JSON(res)
}

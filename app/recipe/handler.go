package recipe

import "github.com/gofiber/fiber/v2"

type RecipeHandler interface {
	AddRecipe(ctx *fiber.Ctx) error
	ListRecipe(ctx *fiber.Ctx) error
	DetailRecipe(ctx *fiber.Ctx) error
	VideoDetailRecipe(ctx *fiber.Ctx) error
	ActionDetailRecipe(ctx *fiber.Ctx) error
	UpdateRecipe(ctx *fiber.Ctx) error
	DeleteRecipe(ctx *fiber.Ctx) error
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

func (c *recipeHandler) DetailRecipe(ctx *fiber.Ctx) error {
	res := c.UseCase.DetailRecipe(ctx)
	return ctx.Status(res.Code).JSON(res)
}

func (c *recipeHandler) VideoDetailRecipe(ctx *fiber.Ctx) error {
	res := c.UseCase.VideoDetailRecipe(ctx)
	return ctx.Status(res.Code).JSON(res)
}

func (c *recipeHandler) ActionDetailRecipe(ctx *fiber.Ctx) error {
	res := c.UseCase.ActionDetailRecipe(ctx)
	return ctx.Status(res.Code).JSON(res)
}

func (c *recipeHandler) UpdateRecipe(ctx *fiber.Ctx) error {
	res := c.UseCase.UpdateRecipe(ctx)
	return ctx.Status(res.Code).JSON(res)
}

func (c *recipeHandler) DeleteRecipe(ctx *fiber.Ctx) error {
	res := c.UseCase.DeleteRecipe(ctx)
	return ctx.Status(res.Code).JSON(res)
}

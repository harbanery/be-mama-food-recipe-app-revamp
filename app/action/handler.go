package action

import "github.com/gofiber/fiber/v2"

type ActionHandler interface {
	ActionSave(ctx *fiber.Ctx) error
	ActionLike(ctx *fiber.Ctx) error
	AddComment(ctx *fiber.Ctx) error
	RemoveComment(ctx *fiber.Ctx) error
}

type actionHandler struct {
	UseCase ActionUseCase
}

func NewActionHandler(useCase ActionUseCase) ActionHandler {
	return &actionHandler{
		UseCase: useCase,
	}
}

func (c *actionHandler) ActionSave(ctx *fiber.Ctx) error {
	res := c.UseCase.ActionSave(ctx)
	return ctx.Status(res.Code).JSON(res)
}

func (c *actionHandler) ActionLike(ctx *fiber.Ctx) error {
	res := c.UseCase.ActionLike(ctx)
	return ctx.Status(res.Code).JSON(res)
}

func (c *actionHandler) AddComment(ctx *fiber.Ctx) error {
	res := c.UseCase.AddComment(ctx)
	return ctx.Status(res.Code).JSON(res)
}

func (c *actionHandler) RemoveComment(ctx *fiber.Ctx) error {
	res := c.UseCase.RemoveComment(ctx)
	return ctx.Status(res.Code).JSON(res)
}

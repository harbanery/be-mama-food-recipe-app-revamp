package account

import "github.com/gofiber/fiber/v2"

type AccountHandler interface {
	DetailProfile(ctx *fiber.Ctx) error
	UpdateProfile(ctx *fiber.Ctx) error
	UpdateProfilePhoto(ctx *fiber.Ctx) error
}

type accountHandler struct {
	UseCase AccountUseCase
}

func NewAccountHandler(useCase AccountUseCase) AccountHandler {
	return &accountHandler{
		UseCase: useCase,
	}
}

func (c *accountHandler) DetailProfile(ctx *fiber.Ctx) error {
	res := c.UseCase.DetailProfile(ctx)
	return ctx.Status(res.Code).JSON(res)
}

func (c *accountHandler) UpdateProfile(ctx *fiber.Ctx) error {
	res := c.UseCase.UpdateProfile(ctx)
	return ctx.Status(res.Code).JSON(res)
}

func (c *accountHandler) UpdateProfilePhoto(ctx *fiber.Ctx) error {
	res := c.UseCase.UpdateProfilePhoto(ctx)
	return ctx.Status(res.Code).JSON(res)
}

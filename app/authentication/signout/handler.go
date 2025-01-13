package signout

import "github.com/gofiber/fiber/v2"

type SignOutHandler interface {
	Logout(ctx *fiber.Ctx) error
}

type signOutHandler struct {
	UseCase SignOutUseCase
}

func NewSignOutHandler(useCase SignOutUseCase) SignOutHandler {
	return &signOutHandler{
		UseCase: useCase,
	}
}

func (c *signOutHandler) Logout(ctx *fiber.Ctx) error {
	res := c.UseCase.Logout(ctx)
	return ctx.Status(res.Code).JSON(res)
}

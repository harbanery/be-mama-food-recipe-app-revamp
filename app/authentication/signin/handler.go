package signin

import "github.com/gofiber/fiber/v2"

type SignInHandler interface {
	Login(ctx *fiber.Ctx) error
}

type signInHandler struct {
	UseCase SignInUseCase
}

func NewSignInHandler(useCase SignInUseCase) SignInHandler {
	return &signInHandler{
		UseCase: useCase,
	}
}

func (c *signInHandler) Login(ctx *fiber.Ctx) error {
	res := c.UseCase.Login(ctx)
	return ctx.Status(res.Code).JSON(res)
}
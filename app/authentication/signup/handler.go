package signup

import "github.com/gofiber/fiber/v2"

type SignUpHandler interface {
	Register(ctx *fiber.Ctx) error
}

type signUpHandler struct {
	UseCase SignUpUseCase
}

func NewSignUpHandler(useCase SignUpUseCase) SignUpHandler {
	return &signUpHandler{
		UseCase: useCase,
	}
}

func (c *signUpHandler) Register(ctx *fiber.Ctx) error {
	res := c.UseCase.Register(ctx)
	return ctx.Status(res.Code).JSON(res)
}
package authentication

import "github.com/gofiber/fiber/v2"

type AuthHandler interface {
	Register(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	Logout(ctx *fiber.Ctx) error
}

type authHandler struct {
	UseCase AuthUseCase
}

func NewAuthHandler(useCase AuthUseCase) AuthHandler {
	return &authHandler{
		UseCase: useCase,
	}
}

func (c *authHandler) Register(ctx *fiber.Ctx) error {
	res := c.UseCase.Register(ctx)
	return ctx.Status(res.Code).JSON(res)
}

func (c *authHandler) Login(ctx *fiber.Ctx) error {
	res := c.UseCase.Login(ctx)
	return ctx.Status(res.Code).JSON(res)
}

func (c *authHandler) Logout(ctx *fiber.Ctx) error {
	res := c.UseCase.Logout(ctx)
	return ctx.Status(res.Code).JSON(res)
}

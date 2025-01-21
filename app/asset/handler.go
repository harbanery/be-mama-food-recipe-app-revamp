package asset

import "github.com/gofiber/fiber/v2"

type AssetHandler interface {
	UploadFile(ctx *fiber.Ctx) error
	RemoveFile(ctx *fiber.Ctx) error
}

type assetHandler struct {
	UseCase AssetUseCase
}

func NewAssetHandler(useCase AssetUseCase) AssetHandler {
	return &assetHandler{
		UseCase: useCase,
	}
}

func (c *assetHandler) UploadFile(ctx *fiber.Ctx) error {
	res := c.UseCase.UploadFile(ctx)
	return ctx.Status(res.Code).JSON(res)
}

func (c *assetHandler) RemoveFile(ctx *fiber.Ctx) error {
	res := c.UseCase.RemoveFile(ctx)
	return ctx.Status(res.Code).JSON(res)
}

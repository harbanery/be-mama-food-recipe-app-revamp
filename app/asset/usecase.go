package asset

import (
	"context"
	"fmt"
	"mama-recipe/helper"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AssetUseCase interface {
	UploadFile(ctx *fiber.Ctx) *helper.WebResponse[interface{}]
	RemoveFile(ctx *fiber.Ctx) *helper.WebResponse[interface{}]
}

type assetUseCase struct {
	DB              *gorm.DB
	Validate        *validator.Validate
	AssetRepository AssetRepository
	Environment     *helper.EnvLoad
	Cloudinary      *cloudinary.Cloudinary
}

func NewAssetUseCase(db *gorm.DB, validate *validator.Validate, assetRepository AssetRepository, env *helper.EnvLoad, cloudinary *cloudinary.Cloudinary) AssetUseCase {
	return &assetUseCase{
		DB:              db,
		Validate:        validate,
		AssetRepository: assetRepository,
		Environment:     env,
		Cloudinary:      cloudinary,
	}
}

type ctxKey string

const requestCtxKey ctxKey = "fasthttp.RequestCtx"

func (c *assetUseCase) UploadFile(ctx *fiber.Ctx) *helper.WebResponse[interface{}] {
	context := context.WithValue(ctx.Context(), requestCtxKey, ctx.Context())

	request := new(AssetRequest)
	if err := helper.ValidateFormRequest(ctx, c.Validate, request); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	files := []*AssetResponse{}
	for _, file := range request.File {
		fileType, err := helper.ValidateFileRequest(file)
		if err != nil {
			return helper.Response(ctx, 400, err.Error(), nil)
		}

		fileUrl, err := helper.UploadFile(&context, c.Cloudinary, file)
		if err != nil {
			return helper.Response(ctx, 400, err.Error(), nil)
		}

		files = append(files, &AssetResponse{
			FileURL:  fileUrl,
			FileType: fileType,
		})
	}

	return helper.Response(ctx, 200, "upload success", files)
}

func (c *assetUseCase) RemoveFile(ctx *fiber.Ctx) *helper.WebResponse[interface{}] {
	context := context.WithValue(ctx.Context(), requestCtxKey, ctx.Context())

	request := new(DeleteAssetRequest)
	if err := helper.ValidateRequest(ctx, c.Validate, request); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	messages := []string{}
	for i, url := range request.FileURL {
		if err := helper.DeleteFile(&context, c.Cloudinary, &url); err != nil {
			messages = append(messages, fmt.Sprint("delete upload no.", i, " failed"))
		} else {
			messages = append(messages, fmt.Sprint("delete upload no.", i, " success"))
		}
	}

	return helper.Response(ctx, 200, "delete upload success", messages)
}

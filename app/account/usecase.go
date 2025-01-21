package account

import (
	"context"
	"fmt"
	"mama-recipe/helper"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AccountUseCase interface {
	DetailProfile(ctx *fiber.Ctx) *helper.WebResponse[interface{}]
	UpdateProfile(ctx *fiber.Ctx) *helper.WebResponse[interface{}]
	UpdateProfilePhoto(ctx *fiber.Ctx) *helper.WebResponse[interface{}]
	DeleteProfilePhoto(ctx *fiber.Ctx) *helper.WebResponse[interface{}]
}

type accountUseCase struct {
	DB                *gorm.DB
	Validate          *validator.Validate
	AccountRepository AccountRepository
	Environment       *helper.EnvLoad
	Cloudinary        *cloudinary.Cloudinary
}

func NewAccountUseCase(db *gorm.DB, validate *validator.Validate, accountRepository AccountRepository, env *helper.EnvLoad, cloudinary *cloudinary.Cloudinary) AccountUseCase {
	return &accountUseCase{
		DB:                db,
		Validate:          validate,
		AccountRepository: accountRepository,
		Environment:       env,
		Cloudinary:        cloudinary,
	}
}

type ctxKey string

const requestCtxKey ctxKey = "fasthttp.RequestCtx"
const profileNotFound string = `profile not found`

func (c *accountUseCase) DetailProfile(ctx *fiber.Ctx) *helper.WebResponse[interface{}] {
	db := c.DB.WithContext(ctx.Context())

	userID := ctx.Locals("id").(string)
	if err := c.AccountRepository.CheckUserID(db, userID); err != nil {
		return helper.Response(ctx, 400, profileNotFound, nil)
	}
	user := c.AccountRepository.GetProfile(db, userID)

	res := helper.Response(ctx, 200, "profile success", &user)
	return res
}

func (c *accountUseCase) UpdateProfile(ctx *fiber.Ctx) *helper.WebResponse[interface{}] {
	db := c.DB.WithContext(ctx.Context())

	userID := ctx.Locals("id").(string)
	if err := c.AccountRepository.CheckUserID(db, userID); err != nil {
		return helper.Response(ctx, 400, profileNotFound, nil)
	}

	request := new(ProfileRequest)
	if err := helper.ValidateRequest(ctx, c.Validate, request); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	if err := c.AccountRepository.UpdateUser(db, userID, *request); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	return helper.Response(ctx, 200, "update profile success", nil)
}

func (c *accountUseCase) UpdateProfilePhoto(ctx *fiber.Ctx) *helper.WebResponse[interface{}] {
	db := c.DB.WithContext(ctx.Context())
	context := context.WithValue(ctx.Context(), requestCtxKey, ctx.Context())

	userID := ctx.Locals("id").(string)
	if err := c.AccountRepository.CheckUserID(db, userID); err != nil {
		return helper.Response(ctx, 400, profileNotFound, nil)
	}

	request := new(PhotoRequest)
	if err := helper.ValidateFormRequest(ctx, c.Validate, request); err != nil {
		fmt.Println("validation photo error")
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	photo := request.Photo[0]
	if _, err := helper.ValidateFileRequest(photo); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	photoUrl, err := helper.UploadFile(&context, c.Cloudinary, photo)
	if err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	if err := c.AccountRepository.UpdateUserPhoto(db, userID, *photoUrl); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	return helper.Response(ctx, 200, "update profile photo success", nil)
}

func (c *accountUseCase) DeleteProfilePhoto(ctx *fiber.Ctx) *helper.WebResponse[interface{}] {
	db := c.DB.WithContext(ctx.Context())
	context := context.WithValue(ctx.Context(), requestCtxKey, ctx.Context())

	userID := ctx.Locals("id").(string)
	if err := c.AccountRepository.CheckUserID(db, userID); err != nil {
		return helper.Response(ctx, 400, profileNotFound, nil)
	}
	user := c.AccountRepository.GetProfile(db, userID)

	if err := helper.DeleteFile(&context, c.Cloudinary, &user.Photo); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	if err := c.AccountRepository.DeleteUserPhoto(db, userID); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	return helper.Response(ctx, 200, "delete profile photo success", nil)
}

package account

import (
	"mama-recipe/helper"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AccountUseCase interface {
	DetailProfile(ctx *fiber.Ctx) *helper.WebResponse[interface{}]
	UpdateProfile(ctx *fiber.Ctx) *helper.WebResponse[interface{}]
}

type accountUseCase struct {
	DB                *gorm.DB
	Validate          *validator.Validate
	AccountRepository AccountRepository
	Environment       *helper.EnvLoad
}

func NewAccountUseCase(db *gorm.DB, validate *validator.Validate, accountRepository AccountRepository, env *helper.EnvLoad) AccountUseCase {
	return &accountUseCase{
		DB:                db,
		Validate:          validate,
		AccountRepository: accountRepository,
		Environment:       env,
	}
}

func (c *accountUseCase) DetailProfile(ctx *fiber.Ctx) *helper.WebResponse[interface{}] {
	db := c.DB.WithContext(ctx.Context())

	userID := ctx.Locals("id").(string)
	user, err := c.AccountRepository.CheckUserFromID(db, userID)
	if err != nil {
		return helper.Response(ctx, 400, "profile not found", nil)
	}

	res := helper.Response(ctx, 200, "profile success", &Profile{
		ID:       user.ID,
		Fullname: user.Fullname,
		Photo:    user.Photo,
	})
	return res
}

func (c *accountUseCase) UpdateProfile(ctx *fiber.Ctx) *helper.WebResponse[interface{}] {
	db := c.DB.WithContext(ctx.Context())

	userID := ctx.Locals("id").(string)
	_, err := c.AccountRepository.CheckUserFromID(db, userID)
	if err != nil {
		return helper.Response(ctx, 400, "profile not found", nil)
	}

	request := new(ProfileRequest)
	err = helper.ValidateRequest(ctx, c.Validate, request)
	if err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	if err := c.AccountRepository.UpdateUserFromID(db, userID, *request); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	res := helper.Response(ctx, 200, "update profile success", nil)
	return res
}

package signout

import (
	"mama-recipe/helper"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SignOutUseCase interface {
	Logout(ctx *fiber.Ctx) *helper.WebResponse[interface{}]
}

type signOutUseCase struct {
	DB                *gorm.DB
	Validate          *validator.Validate
	SignOutRepository SignOutRepository
	Environment       *helper.EnvLoad
}

func NewSignOutUseCase(db *gorm.DB, validate *validator.Validate, signOutRepository SignOutRepository, env *helper.EnvLoad) SignOutUseCase {
	return &signOutUseCase{
		DB: db,
		Validate: validate,
		SignOutRepository: signOutRepository,
		Environment: env,
	}
}

func (c *signOutUseCase) Logout(ctx *fiber.Ctx) *helper.WebResponse[interface{}] {
	response := helper.Response(ctx, 0, "signout success", nil)
	
	return response
}
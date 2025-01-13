package signup

import (
	"mama-recipe/helper"
	"mama-recipe/schema"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SignUpUseCase interface {
	Register(ctx *fiber.Ctx) *helper.WebResponse[interface{}]
}

type signUpUseCase struct {
	DB *gorm.DB
	Validate *validator.Validate
	SignUpRepository SignUpRepository
	Environment              *helper.EnvLoad
}

func NewSignUpUseCase(db *gorm.DB, validate *validator.Validate, signUpRepository SignUpRepository, env *helper.EnvLoad) SignUpUseCase {
	return &signUpUseCase{
		DB: db,
		Validate: validate,
		SignUpRepository: signUpRepository,
		Environment: env,
	}
}

func (c *signUpUseCase) Register(ctx *fiber.Ctx) *helper.WebResponse[interface{}] {
	db := c.DB.WithContext(ctx.Context())

	request := new(SignUp)
	err := helper.ValidateRequest(ctx, c.Validate, request)
	if err != nil {
		return helper.Response(ctx, 400, err.Error(), helper.EmptyObject())
	}

	request.Email = strings.ToLower(request.Email)
	if err := c.SignUpRepository.EmailValidation(db, &request.Email); err != nil {
		return helper.Response(ctx, 400, "email already exist", helper.EmptyObject())
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(request.Password), 12)
	request.Password = "***"

	user := &schema.User{
		Email:    request.Email,
		Password: string(password),
		Username: request.Username,
		Phone:    request.Phone,
	}

	if err := c.SignUpRepository.CreateUser(db, user); err != nil {
		return helper.Response(ctx, 400, err.Error(), helper.EmptyObject())
	}

	res := helper.Response(ctx, 200, "signup success", nil)
	return res
}

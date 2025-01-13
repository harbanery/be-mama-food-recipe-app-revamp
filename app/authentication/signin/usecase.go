package signin

import (
	"mama-recipe/helper"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SignInUseCase interface {
	Login(ctx *fiber.Ctx) *helper.WebResponse[interface{}]
}

type signInUseCase struct {
	DB               *gorm.DB
	Validate         *validator.Validate
	SignInRepository SignInRepository
	Environment              *helper.EnvLoad
}

func NewSignInUseCase(db *gorm.DB, validate *validator.Validate, signInRepository SignInRepository, env *helper.EnvLoad) SignInUseCase {
	return &signInUseCase{
		DB:               db,
		Validate:         validate,
		SignInRepository: signInRepository,
		Environment: env,
	}
}

func (c *signInUseCase) Login(ctx *fiber.Ctx) *helper.WebResponse[interface{}] {
	db := c.DB.WithContext(ctx.Context())

	request := new(SignIn)
	err := helper.ValidateRequest(ctx, c.Validate, request)
	if err != nil {
		return helper.Response(ctx, 400, err.Error(), helper.EmptyObject())
	}

	request.Email = strings.ToLower(request.Email)
	user, err := c.SignInRepository.CheckEmail(db, &request.Email); 
	if err != nil {
		return helper.Response(ctx, 400, "invalid email or password", helper.EmptyObject())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return helper.Response(ctx, 401, "invalid email or password", helper.EmptyObject())
	}

	if c.Environment == nil {
		return helper.Response(ctx, 500, "environment not found", helper.EmptyObject())
	}

	token, _ := helper.GenerateToken(c.Environment.JWT_SECRET_KEY, map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
		"exp": time.Now().Add(time.Hour * 12).Unix(),
	})

	refreshToken, _ := helper.GenerateToken(c.Environment.JWT_SECRET_KEY, map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	res := helper.Response(ctx, 200, "signin success", fiber.Map{
		"email": request.Email,
		"access_token": token,
		"refresh_token": refreshToken,
	})
	return res
}
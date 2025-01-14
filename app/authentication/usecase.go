package authentication

import (
	"mama-recipe/helper"
	"mama-recipe/schema"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUseCase interface {
	Register(ctx *fiber.Ctx) *helper.WebResponse[interface{}]
	Login(ctx *fiber.Ctx) *helper.WebResponse[interface{}]
	Logout(ctx *fiber.Ctx) *helper.WebResponse[interface{}]
}

type authUseCase struct {
	DB             *gorm.DB
	Validate       *validator.Validate
	AuthRepository AuthRepository
	Environment    *helper.EnvLoad
}

func NewAuthUseCase(db *gorm.DB, validate *validator.Validate, authRepository AuthRepository, env *helper.EnvLoad) AuthUseCase {
	return &authUseCase{
		DB:             db,
		Validate:       validate,
		AuthRepository: authRepository,
		Environment:    env,
	}
}

func (c *authUseCase) Register(ctx *fiber.Ctx) *helper.WebResponse[interface{}] {
	db := c.DB.WithContext(ctx.Context())

	request := new(SignUpRequest)
	err := helper.ValidateRequest(ctx, c.Validate, request)
	if err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	request.Email = strings.ToLower(request.Email)
	if err := c.AuthRepository.EmailValidation(db, &request.Email); err != nil {
		return helper.Response(ctx, 400, "email already exist", nil)
	}

	username := helper.UsernameFromEmail(request.Email)
	password, _ := bcrypt.GenerateFromPassword([]byte(request.Password), 12)
	request.Password = "***"

	user := &schema.User{
		Email:    request.Email,
		Password: string(password),
		Username: username,
		Fullname: request.Fullname,
		Phone:    request.Phone,
	}

	if err := c.AuthRepository.CreateUser(db, user); err != nil {
		return helper.Response(ctx, 400, err.Error(), helper.EmptyObject())
	}

	res := helper.Response(ctx, 200, "signup success", nil)
	return res
}

func (c *authUseCase) Login(ctx *fiber.Ctx) *helper.WebResponse[interface{}] {
	db := c.DB.WithContext(ctx.Context())

	request := new(SignInRequest)
	err := helper.ValidateRequest(ctx, c.Validate, request)
	if err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	request.Email = strings.ToLower(request.Email)
	user, err := c.AuthRepository.CheckEmail(db, &request.Email)
	if err != nil {
		return helper.Response(ctx, 400, "invalid email or password", nil)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return helper.Response(ctx, 401, "invalid email or password", nil)
	}

	if c.Environment == nil {
		return helper.Response(ctx, 500, "environment not found", nil)
	}

	tokenExp, err := strconv.Atoi(c.Environment.JWT_TOKEN_EXPIRATION)
	if err != nil {
		return helper.Response(ctx, 500, "environment not found", nil)
	}
	timeTokenExp := time.Now().Add(time.Hour * time.Duration(tokenExp))

	token, err := helper.GenerateToken(c.Environment.JWT_TOKEN_SECRET_KEY, map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
		"exp":   timeTokenExp.Unix(),
	})
	if err != nil {
		return helper.Response(ctx, 500, "token not generated", nil)
	}

	timeFormat := timeTokenExp.Format("2006-01-02 15:04:05")
	res := helper.Response(ctx, 200, "signin success", &SignInResponse{
		Email:                  request.Email,
		AccessToken:            token,
		AccessTokenTimeExpired: &timeFormat,
	})
	return res
}

func (c *authUseCase) Logout(ctx *fiber.Ctx) *helper.WebResponse[interface{}] {
	response := helper.Response(ctx, 0, "signout success", nil)

	return response
}

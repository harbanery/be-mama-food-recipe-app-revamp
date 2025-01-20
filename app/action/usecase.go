package action

import (
	"mama-recipe/helper"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ActionUseCase interface {
	ActionSave(ctx *fiber.Ctx) *helper.WebResponse[interface{}]
	ActionLike(ctx *fiber.Ctx) *helper.WebResponse[interface{}]
	AddComment(ctx *fiber.Ctx) *helper.WebResponse[interface{}]
	RemoveComment(ctx *fiber.Ctx) *helper.WebResponse[interface{}]
}

type actionUseCase struct {
	DB               *gorm.DB
	Validate         *validator.Validate
	ActionRepository ActionRepository
	Environment      *helper.EnvLoad
}

func NewActionUseCase(db *gorm.DB, validate *validator.Validate, actionRepository ActionRepository, env *helper.EnvLoad) ActionUseCase {
	return &actionUseCase{
		DB:               db,
		Validate:         validate,
		ActionRepository: actionRepository,
		Environment:      env,
	}
}

func (c *actionUseCase) ActionSave(ctx *fiber.Ctx) *helper.WebResponse[interface{}] {
	db := c.DB.WithContext(ctx.Context())

	userID := ctx.Locals("id").(string)
	if err := c.ActionRepository.CheckUser(db, userID); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	request := new(ActionRequest)
	if err := helper.ValidateRequest(ctx, c.Validate, request); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	if err := c.ActionRepository.CheckRecipe(db, request.RecipeID); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	id, err := c.ActionRepository.CheckSave(db, request.RecipeID, userID)
	if err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	} else if id != nil {
		if err := c.ActionRepository.DeleteSave(db, request.RecipeID, userID); err != nil {
			return helper.Response(ctx, 400, err.Error(), nil)
		}

		return helper.Response(ctx, 200, "save removed success", nil)
	}

	if err := c.ActionRepository.CreateSave(db, request.RecipeID, userID); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	return helper.Response(ctx, 200, "save success", nil)
}

func (c *actionUseCase) ActionLike(ctx *fiber.Ctx) *helper.WebResponse[interface{}] {
	db := c.DB.WithContext(ctx.Context())

	userID := ctx.Locals("id").(string)
	if err := c.ActionRepository.CheckUser(db, userID); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	request := new(ActionRequest)
	if err := helper.ValidateRequest(ctx, c.Validate, request); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	if err := c.ActionRepository.CheckRecipe(db, request.RecipeID); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	id, err := c.ActionRepository.CheckLike(db, request.RecipeID, userID)
	if err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	} else if id != nil {
		if err := c.ActionRepository.DeleteLike(db, request.RecipeID, userID); err != nil {

			return helper.Response(ctx, 400, err.Error(), nil)
		}

		return helper.Response(ctx, 200, "like removed success", nil)
	}

	if err := c.ActionRepository.CreateLike(db, request.RecipeID, userID); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	return helper.Response(ctx, 200, "like success", nil)
}

func (c *actionUseCase) AddComment(ctx *fiber.Ctx) *helper.WebResponse[interface{}] {
	db := c.DB.WithContext(ctx.Context())

	userID := ctx.Locals("id").(string)
	if err := c.ActionRepository.CheckUser(db, userID); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	request := new(CommentRequest)
	if err := helper.ValidateRequest(ctx, c.Validate, request); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	if err := c.ActionRepository.CheckRecipe(db, request.RecipeID); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	if err := c.ActionRepository.CreateComment(db, request, userID); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	return helper.Response(ctx, 200, "comment success", nil)
}

func (c *actionUseCase) RemoveComment(ctx *fiber.Ctx) *helper.WebResponse[interface{}] {
	db := c.DB.WithContext(ctx.Context())

	userID := ctx.Locals("id").(string)
	if err := c.ActionRepository.CheckUser(db, userID); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	request := new(ActionRequest)
	if err := helper.ValidateRequest(ctx, c.Validate, request); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	if err := c.ActionRepository.CheckRecipe(db, request.RecipeID); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	if err := c.ActionRepository.DeleteComment(db, request.RecipeID, userID); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	return helper.Response(ctx, 200, "comment removed success", nil)
}

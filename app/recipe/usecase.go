package recipe

import (
	"context"
	"mama-recipe/helper"
	"mama-recipe/schema"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type RecipeUseCase interface {
	AddRecipe(ctx *fiber.Ctx) *helper.WebResponse[interface{}]
	ListRecipe(ctx *fiber.Ctx) *helper.WebResponse[interface{}]
	DetailRecipe(ctx *fiber.Ctx) *helper.WebResponse[interface{}]
	ActionDetailRecipe(ctx *fiber.Ctx) *helper.WebResponse[interface{}]
	UpdateRecipe(ctx *fiber.Ctx) *helper.WebResponse[interface{}]
	DeleteRecipe(ctx *fiber.Ctx) *helper.WebResponse[interface{}]
}

type recipeUseCase struct {
	DB               *gorm.DB
	Validate         *validator.Validate
	RecipeRepository RecipeRepository
	Environment      *helper.EnvLoad
	Cloudinary       *cloudinary.Cloudinary
}

func NewRecipeUseCase(db *gorm.DB, validate *validator.Validate, recipeRepository RecipeRepository, env *helper.EnvLoad, cloudinary *cloudinary.Cloudinary) RecipeUseCase {
	return &recipeUseCase{
		DB:               db,
		Validate:         validate,
		RecipeRepository: recipeRepository,
		Environment:      env,
		Cloudinary:       cloudinary,
	}
}

type ctxKey string

const requestCtxKey ctxKey = "fasthttp.RequestCtx"

func (c *recipeUseCase) AddRecipe(ctx *fiber.Ctx) *helper.WebResponse[interface{}] {
	db := c.DB.WithContext(ctx.Context())
	context := context.WithValue(ctx.Context(), requestCtxKey, ctx.Context())

	userID := ctx.Locals("id").(string)
	email := ctx.Locals("email").(string)
	if err := c.RecipeRepository.CheckUser(db, userID, &email); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	request := new(RecipeRequest)
	if err := helper.ValidateFormRequest(ctx, c.Validate, request); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	username := helper.UsernameFromEmail(email)
	slug := helper.ToSlug(request.Title, username)
	photo := request.Image[0]
	if err := helper.ValidateImageRequest(photo); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	photoUrl, err := helper.UploadFile(&context, c.Cloudinary, photo)
	if err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	recipe := &schema.Recipe{
		Title:       request.Title,
		SubTitle:    request.SubTitle,
		Slug:        slug,
		Header:      request.Header,
		Image:       *photoUrl,
		Description: request.Description,
		AuthorID:    userID,
	}

	if err := c.RecipeRepository.CreateRecipe(db, recipe); err != nil {
		return helper.Response(ctx, 400, err.Error(), helper.EmptyObject())
	}

	return helper.Response(ctx, 200, "add recipe success", nil)
}

func (c *recipeUseCase) ListRecipe(ctx *fiber.Ctx) *helper.WebResponse[interface{}] {
	db := c.DB.WithContext(ctx.Context())

	recipes, err := c.RecipeRepository.GetRecipes(db)
	if err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	return helper.Response(ctx, 200, "list recipe success", recipes)
}

func (c *recipeUseCase) DetailRecipe(ctx *fiber.Ctx) *helper.WebResponse[interface{}] {
	db := c.DB.WithContext(ctx.Context())
	slug := ctx.Query("slug")
	if slug == "" {
		return helper.Response(ctx, 400, "slug is required", nil)
	}

	recipe, err := c.RecipeRepository.DetailRecipe(db, slug)
	if err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	return helper.Response(ctx, 200, "detail recipe success", recipe)
}

func (c *recipeUseCase) ActionDetailRecipe(ctx *fiber.Ctx) *helper.WebResponse[interface{}] {
	db := c.DB.WithContext(ctx.Context())
	recipeID := ctx.Query("id")
	if recipeID == "" {
		return helper.Response(ctx, 400, "recipe id is required", nil)
	}

	userID := ctx.Locals("id").(string)
	if err := c.RecipeRepository.CheckUser(db, userID); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	var actionResponse = &ActionRecipeResponse{
		IsSaved: false,
		IsLiked: false,
	}

	saveID, err := c.RecipeRepository.CheckSave(db, recipeID, userID)
	if err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	} else if saveID != nil {
		actionResponse.IsLiked = true
	}

	likeID, err := c.RecipeRepository.CheckLike(db, recipeID, userID)
	if err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	} else if likeID != nil {
		actionResponse.IsSaved = true
	}

	return helper.Response(ctx, 200, "action recipe success", actionResponse)
}

func (c *recipeUseCase) UpdateRecipe(ctx *fiber.Ctx) *helper.WebResponse[interface{}] {
	db := c.DB.WithContext(ctx.Context())
	context := context.WithValue(ctx.Context(), requestCtxKey, ctx.Context())

	userID := ctx.Locals("id").(string)
	if err := c.RecipeRepository.CheckUser(db, userID); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	request := new(RecipeRequest)
	if err := helper.ValidateFormRequest(ctx, c.Validate, request); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	recipe, err := c.RecipeRepository.GetRecipe(db, request.ID)
	if err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	if recipe.AuthorID != userID {
		return helper.Response(ctx, 400, "unauthorized", nil)
	}

	photo := request.Image[0]
	if err := helper.ValidateImageRequest(photo); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	photoUrl, err := helper.UploadFile(&context, c.Cloudinary, photo)
	if err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	if photoUrl == &recipe.Image {
		if err := helper.DeleteFile(&context, c.Cloudinary, &recipe.Image); err != nil {
			return helper.Response(ctx, 400, err.Error(), nil)
		}
	}

	recipe.Title = request.Title
	recipe.SubTitle = request.SubTitle
	recipe.Header = request.Header
	recipe.Description = request.Description

	if err := c.RecipeRepository.EditRecipe(db, recipe.ID, recipe); err != nil {
		return helper.Response(ctx, 400, err.Error(), helper.EmptyObject())
	}

	return helper.Response(ctx, 200, "edit recipe success", nil)
}

func (c *recipeUseCase) DeleteRecipe(ctx *fiber.Ctx) *helper.WebResponse[interface{}] {
	db := c.DB.WithContext(ctx.Context())
	context := context.WithValue(ctx.Context(), requestCtxKey, ctx.Context())
	id := ctx.Query("id")
	if id == "" {
		return helper.Response(ctx, 400, "id is required", nil)
	}

	userID := ctx.Locals("id").(string)
	if err := c.RecipeRepository.CheckUser(db, userID); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	recipe, err := c.RecipeRepository.GetRecipe(db, id)
	if err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	if recipe.AuthorID != userID {
		return helper.Response(ctx, 400, "unauthorized", nil)
	}

	if err := helper.DeleteFile(&context, c.Cloudinary, &recipe.Image); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	if err := c.RecipeRepository.DeleteRecipe(db, recipe.ID); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	return helper.Response(ctx, 200, "delete recipe success", nil)
}

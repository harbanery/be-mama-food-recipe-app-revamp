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
	VideoDetailRecipe(ctx *fiber.Ctx) *helper.WebResponse[interface{}]
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

	if len(request.Image) > 1 {
		return helper.Response(ctx, 400, "only 1 image allowed", nil)
	}

	photo := request.Image[0]
	if err := helper.ValidateImageRequest(photo); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	photoUrl, err := helper.UploadFile(&context, c.Cloudinary, photo)
	if err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	videos := []*schema.Video{}
	for _, video := range request.Video {
		if err := helper.ValidateVideoRequest(video); err != nil {
			return helper.Response(ctx, 400, err.Error(), nil)
		}

		videoUrl, err := helper.UploadFile(&context, c.Cloudinary, video)
		if err != nil {
			return helper.Response(ctx, 400, err.Error(), nil)
		}

		videos = append(videos, &schema.Video{
			Title:  request.Title,
			Source: "upload",
			URL:    *videoUrl,
		})
	}

	for _, videoUrl := range request.VideoURL {
		videos = append(videos, &schema.Video{
			Title:  request.Title,
			Source: "url",
			URL:    videoUrl,
		})
	}

	recipe := &schema.Recipe{
		Title:       request.Title,
		SubTitle:    request.SubTitle,
		Slug:        slug,
		Header:      request.Header,
		Image:       *photoUrl,
		Videos:      videos,
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

	paramsrequest := helper.NewParamsRequest(ctx)

	recipes, counts, err := c.RecipeRepository.GetRecipes(db, paramsrequest)
	if err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	response := helper.NewParamsResponse(recipes, counts, paramsrequest)

	return helper.Response(ctx, 200, "list recipe success", response)
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

func (c *recipeUseCase) VideoDetailRecipe(ctx *fiber.Ctx) *helper.WebResponse[interface{}] {
	db := c.DB.WithContext(ctx.Context())
	videoID := ctx.Query("video_id")
	if videoID == "" {
		return helper.Response(ctx, 400, "video id is required", nil)
	}

	recipeID := ctx.Query("recipe_id")
	if recipeID == "" {
		return helper.Response(ctx, 400, "recipe id is required", nil)
	}

	response, err := c.RecipeRepository.DetailVideo(db, videoID, recipeID)
	if err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	return helper.Response(ctx, 200, "detail video recipe success", response)
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
		actionResponse.IsSaved = true
	}

	likeID, err := c.RecipeRepository.CheckLike(db, recipeID, userID)
	if err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	} else if likeID != nil {
		actionResponse.IsLiked = true
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

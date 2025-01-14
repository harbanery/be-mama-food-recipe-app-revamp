package recipe

import (
	"mama-recipe/helper"
	"mama-recipe/schema"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type RecipeUseCase interface {
	AddRecipe(ctx *fiber.Ctx) *helper.WebResponse[interface{}]
	ListRecipe(ctx *fiber.Ctx) *helper.WebResponse[interface{}]
}

type recipeUseCase struct {
	DB               *gorm.DB
	Validate         *validator.Validate
	RecipeRepository RecipeRepository
	Environment      *helper.EnvLoad
}

func NewRecipeUseCase(db *gorm.DB, validate *validator.Validate, recipeRepository RecipeRepository, env *helper.EnvLoad) RecipeUseCase {
	return &recipeUseCase{
		DB:               db,
		Validate:         validate,
		RecipeRepository: recipeRepository,
		Environment:      env,
	}
}

func (c *recipeUseCase) AddRecipe(ctx *fiber.Ctx) *helper.WebResponse[interface{}] {
	db := c.DB.WithContext(ctx.Context())

	userID := ctx.Locals("id").(string)
	email := ctx.Locals("email").(string)
	if err := c.RecipeRepository.CheckUser(db, userID, email); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	request := new(RecipeRequest)
	if err := helper.ValidateRequest(ctx, c.Validate, request); err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	username := helper.UsernameFromEmail(email)
	slug := helper.ToSlug(request.Title, username)

	recipe := &schema.Recipe{
		Title:       request.Title,
		SubTitle:    request.SubTitle,
		Slug:        slug,
		Header:      request.Header,
		Image:       request.Image,
		Description: request.Description,
		AuthorID:    userID,
	}

	if err := c.RecipeRepository.CreateRecipe(db, recipe); err != nil {
		return helper.Response(ctx, 400, err.Error(), helper.EmptyObject())
	}

	res := helper.Response(ctx, 200, "add recipe success", nil)
	return res
}

func (c *recipeUseCase) ListRecipe(ctx *fiber.Ctx) *helper.WebResponse[interface{}] {
	db := c.DB.WithContext(ctx.Context())

	recipes, err := c.RecipeRepository.GetRecipes(db)
	if err != nil {
		return helper.Response(ctx, 400, err.Error(), nil)
	}

	var recipeResponses []RecipeResponse
	for _, recipe := range recipes {
		recipeResponse := RecipeResponse{
			ID:     recipe.ID,
			Title:  recipe.Title,
			Slug:   recipe.Slug,
			Image:  recipe.Image,
			Author: recipe.Author.Fullname,
		}
		recipeResponses = append(recipeResponses, recipeResponse)
	}

	return helper.Response(ctx, 200, "list recipe success", recipeResponses)
}

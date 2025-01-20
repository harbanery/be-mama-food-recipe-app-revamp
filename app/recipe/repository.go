package recipe

import (
	"mama-recipe/helper"
	"mama-recipe/schema"

	"gorm.io/gorm"
)

type RecipeRepository interface {
	CheckUser(db *gorm.DB, id string, email ...*string) error
	GetRecipes(db *gorm.DB, params *helper.ParamsRequest) ([]*RecipeResponse, *helper.RecordCount, error)
	GetRecipe(db *gorm.DB, id string) (*schema.Recipe, error)
	DetailRecipe(db *gorm.DB, slug string) (*DetailRecipeResponse, error)
	CreateRecipe(db *gorm.DB, data *schema.Recipe) error
	EditRecipe(db *gorm.DB, id string, data *schema.Recipe) error
	DeleteRecipe(db *gorm.DB, id string) error
	CheckSave(db *gorm.DB, recipeID, userID string) (*string, error)
	CheckLike(db *gorm.DB, recipeID, userID string) (*string, error)
}

type recipeRepository struct{}

func NewRecipeRepository() RecipeRepository {
	return &recipeRepository{}
}

const idQueryParam string = "id = ?"
const idActionQueryParams string = "recipe_id = ? AND user_id = ?"

func (s *recipeRepository) CheckUser(db *gorm.DB, id string, email ...*string) error {
	var user *schema.User
	db.First(&user, idQueryParam, id)

	if email != nil {
		db.Where("email = ?", email)
	}

	return db.Error
}

func (s *recipeRepository) GetRecipes(db *gorm.DB, params *helper.ParamsRequest) ([]*RecipeResponse, *helper.RecordCount, error) {
	var recipes []*RecipeResponse
	var recipeCounts helper.RecordCount

	query := db.Model(&schema.Recipe{})

	if params.Sort == "popularity" {
		query = query.Order("(SELECT COUNT(*) FROM saves WHERE saves.recipe_id = recipes.id) " + params.OrderBy).
			Order("(SELECT COUNT(*) FROM likes WHERE likes.recipe_id = recipes.id) " + params.OrderBy)
	} else {
		query = query.Order(params.Sort + " " + params.OrderBy)
	}

	query = query.Offset((params.Page - 1) * params.Limit).Limit(params.Limit).Find(&recipes).Count(&recipeCounts.TotalData)
	recipeCounts.FilteredData = int64(len(recipes))

	if query.Error != nil {
		return nil, nil, query.Error
	}

	return recipes, &recipeCounts, nil
}

func (s *recipeRepository) GetRecipe(db *gorm.DB, id string) (*schema.Recipe, error) {
	var recipe *schema.Recipe
	db.First(&recipe, idQueryParam, id)

	if recipe.ID == "" {
		return nil, gorm.ErrRecordNotFound
	}

	return recipe, nil
}

func (s *recipeRepository) DetailRecipe(db *gorm.DB, slug string) (*DetailRecipeResponse, error) {
	var recipe *DetailRecipeResponse
	db.Model(&schema.Recipe{}).Preload("Author", func(db *gorm.DB) *gorm.DB {
		var author []*Author
		return db.Model(&schema.User{}).Find(&author)
	}).Preload("Saves").Preload("Likes").Preload("Comments.User", func(db *gorm.DB) *gorm.DB {
		var user []*User
		return db.Model(&schema.User{}).Find(&user)
	}).First(&recipe, "slug = ?", slug)

	if recipe.ID == "" {
		return nil, gorm.ErrRecordNotFound
	}

	return recipe, nil
}

func (s *recipeRepository) CreateRecipe(db *gorm.DB, data *schema.Recipe) error {
	return db.Create(&data).Error
}

func (s *recipeRepository) EditRecipe(db *gorm.DB, id string, data *schema.Recipe) error {
	return db.Model(&schema.Recipe{}).Where(idQueryParam, id).Updates(data).Error
}

func (s *recipeRepository) DeleteRecipe(db *gorm.DB, id string) error {
	return db.Delete(&schema.Recipe{}, idQueryParam, id).Error
}

func (s *recipeRepository) CheckSave(db *gorm.DB, recipeID, userID string) (*string, error) {
	var save *schema.Save
	db.First(&save, idActionQueryParams, recipeID, userID)

	if save.ID == "" || db.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &save.ID, nil
}

func (s *recipeRepository) CheckLike(db *gorm.DB, recipeID, userID string) (*string, error) {
	var save *schema.Like
	db.First(&save, idActionQueryParams, recipeID, userID)

	if save.ID == "" || db.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &save.ID, nil
}

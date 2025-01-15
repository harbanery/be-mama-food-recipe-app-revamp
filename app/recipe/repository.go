package recipe

import (
	"mama-recipe/schema"

	"gorm.io/gorm"
)

type RecipeRepository interface {
	CheckUser(db *gorm.DB, id string, email ...*string) error
	GetRecipes(db *gorm.DB) ([]*RecipeResponse, error)
	GetRecipe(db *gorm.DB, id string) (*schema.Recipe, error)
	DetailRecipe(db *gorm.DB, slug string) (*DetailRecipeResponse, error)
	CreateRecipe(db *gorm.DB, data *schema.Recipe) error
	EditRecipe(db *gorm.DB, id string, data *schema.Recipe) error
	DeleteRecipe(db *gorm.DB, id string) error
}

type recipeRepository struct{}

func NewRecipeRepository() RecipeRepository {
	return &recipeRepository{}
}

const idQueryParams string = "id = ?"

func (s *recipeRepository) CheckUser(db *gorm.DB, id string, email ...*string) error {
	var user *schema.User
	db.First(&user, idQueryParams, id)

	if email != nil {
		db.Where("email = ?", email)
	}

	return db.Error
}

func (s *recipeRepository) GetRecipes(db *gorm.DB) ([]*RecipeResponse, error) {
	var recipes []*RecipeResponse
	db.Model(&schema.Recipe{}).Preload("Author", func(db *gorm.DB) *gorm.DB {
		var author []*Author
		return db.Model(&schema.User{}).Find(&author)
	}).Order("created_at desc").Find(&recipes)
	return recipes, nil
}

func (s *recipeRepository) GetRecipe(db *gorm.DB, id string) (*schema.Recipe, error) {
	var recipe *schema.Recipe
	db.First(&recipe, idQueryParams, id)

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
	return db.Model(&schema.Recipe{}).Where(idQueryParams, id).Updates(data).Error
}

func (s *recipeRepository) DeleteRecipe(db *gorm.DB, id string) error {
	return db.Delete(&schema.Recipe{}, "id = ?", id).Error
}

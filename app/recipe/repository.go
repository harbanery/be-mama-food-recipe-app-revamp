package recipe

import (
	"mama-recipe/schema"

	"gorm.io/gorm"
)

type RecipeRepository interface {
	CheckUser(db *gorm.DB, id, email string) error
	CreateRecipe(db *gorm.DB, recipe *schema.Recipe) error
	GetRecipes(db *gorm.DB) ([]*schema.Recipe, error)
}

type recipeRepository struct{}

func NewRecipeRepository() RecipeRepository {
	return &recipeRepository{}
}

func (s *recipeRepository) CheckUser(db *gorm.DB, id, email string) error {
	var user *schema.User
	return db.First(&user, "id = ? AND email = ?", id, email).Error
}

func (s *recipeRepository) CreateRecipe(db *gorm.DB, recipe *schema.Recipe) error {
	return db.Create(&recipe).Error
}

func (s *recipeRepository) GetRecipes(db *gorm.DB) ([]*schema.Recipe, error) {
	var recipes []*schema.Recipe
	db.Preload("Author", func(db *gorm.DB) *gorm.DB {
		var author []*Author
		return db.Model(&schema.User{}).Find(&author)
	}).Find(&recipes)
	return recipes, nil
}

func (s *recipeRepository) DetailRecipes(db *gorm.DB, slug string) (*schema.Recipe, error) {
	var recipe *schema.Recipe
	db.First(&recipe, "slug = ?", slug)
	return recipe, nil
}

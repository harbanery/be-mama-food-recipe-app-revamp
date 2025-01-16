package action

import (
	"mama-recipe/schema"

	"gorm.io/gorm"
)

type ActionRepository interface {
	CheckUser(db *gorm.DB, id string) error
	CheckRecipe(db *gorm.DB, id string) error
	CheckSave(db *gorm.DB, recipeID, userID string) (*string, error)
	CreateSave(db *gorm.DB, recipeID, userID string) error
	DeleteSave(db *gorm.DB, recipeID, userID string) error
	CheckLike(db *gorm.DB, recipeID, userID string) (*string, error)
	CreateLike(db *gorm.DB, recipeID, userID string) error
	DeleteLike(db *gorm.DB, recipeID, userID string) error
}

type actionRepository struct{}

func NewActionRepository() ActionRepository {
	return &actionRepository{}
}

const idQueryParam string = "id = ?"
const idActionQueryParams string = "recipe_id = ? AND user_id = ?"

func (s *actionRepository) CheckUser(db *gorm.DB, id string) error {
	var user *schema.User
	return db.First(&user, idQueryParam, id).Error
}

func (s *actionRepository) CheckRecipe(db *gorm.DB, id string) error {
	var recipe *schema.Recipe
	db.First(&recipe, idQueryParam, id)

	if recipe.ID == "" {
		return gorm.ErrRecordNotFound
	}

	return db.Error
}

func (s *actionRepository) CheckSave(db *gorm.DB, recipeID, userID string) (*string, error) {
	var save *schema.Save
	db.First(&save, idActionQueryParams, recipeID, userID)

	if save.ID == "" || db.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &save.ID, nil
}

func (s *actionRepository) CreateSave(db *gorm.DB, recipeID, userID string) error {
	data := &schema.Save{
		RecipeID: recipeID,
		UserID:   userID,
	}

	return db.Create(&data).Error
}

func (s *actionRepository) DeleteSave(db *gorm.DB, recipeID, userID string) error {
	return db.Unscoped().Delete(&schema.Save{}, idActionQueryParams, recipeID, userID).Error
}

func (s *actionRepository) CheckLike(db *gorm.DB, recipeID, userID string) (*string, error) {
	var save *schema.Like
	db.First(&save, idActionQueryParams, recipeID, userID)

	if save.ID == "" || db.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &save.ID, nil
}

func (s *actionRepository) CreateLike(db *gorm.DB, recipeID, userID string) error {
	data := &schema.Like{
		RecipeID: recipeID,
		UserID:   userID,
	}

	return db.Create(&data).Error
}

func (s *actionRepository) DeleteLike(db *gorm.DB, recipeID, userID string) error {
	return db.Unscoped().Delete(&schema.Like{}, idActionQueryParams, recipeID, userID).Error
}

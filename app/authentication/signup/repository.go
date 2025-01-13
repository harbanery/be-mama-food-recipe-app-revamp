package signup

import (
	"mama-recipe/schema"

	"gorm.io/gorm"
)

type SignUpRepository interface {
	EmailValidation(db *gorm.DB, email *string) error
	CreateUser(db *gorm.DB, user *schema.User) error
}

type signUpRepository struct{}

func NewSignUpRepository() SignUpRepository {
	return &signUpRepository{}
}

func (s *signUpRepository) EmailValidation(db *gorm.DB, email *string) error {
	var user *schema.User
	db.First(&user, "email = ?", email)
	if user.Email != "" {
		return gorm.ErrRegistered
	}
	return nil
}

func (s *signUpRepository) CreateUser(db *gorm.DB, user *schema.User) error {
	result := db.Create(&user)
	return result.Error
}

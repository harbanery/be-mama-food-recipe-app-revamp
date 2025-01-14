package authentication

import (
	"mama-recipe/schema"

	"gorm.io/gorm"
)

type AuthRepository interface {
	EmailValidation(db *gorm.DB, email *string) error
	CreateUser(db *gorm.DB, user *schema.User) error
	CheckEmail(db *gorm.DB, email *string) (schema.User, error)
}

type authRepository struct{}

func NewAuthRepository() AuthRepository {
	return &authRepository{}
}

func (s *authRepository) EmailValidation(db *gorm.DB, email *string) error {
	var user *schema.User
	db.First(&user, "email = ?", email)
	if user.Email != "" {
		return gorm.ErrRegistered
	}
	return nil
}

func (s *authRepository) CheckEmail(db *gorm.DB, email *string) (schema.User, error) {
	var user *schema.User
	db.First(&user, "email = ?", email)
	if user.Email == "" {
		return schema.User{}, gorm.ErrRecordNotFound
	}
	return *user, nil
}

func (s *authRepository) CreateUser(db *gorm.DB, user *schema.User) error {
	result := db.Create(&user)
	return result.Error
}

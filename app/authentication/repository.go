package authentication

import (
	"mama-recipe/schema"

	"gorm.io/gorm"
)

type AuthRepository interface {
	EmailValidation(db *gorm.DB, email *string) error
	CheckEmail(db *gorm.DB, email *string) error
	GetUserFromEmail(db *gorm.DB, email *string) *User
	CreateUser(db *gorm.DB, user *schema.User) error
}

type authRepository struct{}

func NewAuthRepository() AuthRepository {
	return &authRepository{}
}

const emailQueryParam = `email = ?`

func (s *authRepository) EmailValidation(db *gorm.DB, email *string) error {
	var user *schema.User
	db.First(&user, emailQueryParam, email)
	if user.Email != "" {
		return gorm.ErrRegistered
	}
	return nil
}

func (s *authRepository) CheckEmail(db *gorm.DB, email *string) error {
	var user *schema.User
	return db.First(&user, emailQueryParam, email).Error
}

func (s *authRepository) GetUserFromEmail(db *gorm.DB, email *string) *User {
	var user *User
	db.Model(&schema.User{}).First(&user, emailQueryParam, email)
	return user
}

func (s *authRepository) CreateUser(db *gorm.DB, user *schema.User) error {
	result := db.Create(&user)
	return result.Error
}

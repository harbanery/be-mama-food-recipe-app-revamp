package signin

import (
	"mama-recipe/schema"

	"gorm.io/gorm"
)

type SignInRepository interface {
	CheckEmail(db *gorm.DB, email *string) (schema.User, error)
}

type signInRepository struct{}

func NewSignInRepository() SignInRepository {
	return &signInRepository{}
}

func (s *signInRepository) CheckEmail(db *gorm.DB, email *string) (schema.User, error) {
	var user *schema.User
	db.First(&user, "email = ?", email)
	if user.Email == "" {
		return schema.User{}, gorm.ErrRecordNotFound
	}
	return *user, nil
}

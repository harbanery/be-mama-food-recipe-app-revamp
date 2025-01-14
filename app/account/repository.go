package account

import (
	"mama-recipe/schema"

	"gorm.io/gorm"
)

type AccountRepository interface {
	CheckUserFromID(db *gorm.DB, id string) (schema.User, error)
	UpdateUserFromID(db *gorm.DB, id string, data ProfileRequest) error
}

type accountRepository struct{}

func NewAccountRepository() AccountRepository {
	return &accountRepository{}
}

func (s *accountRepository) CheckUserFromID(db *gorm.DB, id string) (schema.User, error) {
	var user *schema.User
	db.First(&user, "id = ?", id)
	if user.ID == "" {
		return schema.User{}, gorm.ErrRecordNotFound
	}

	return *user, nil
}

func (s *accountRepository) UpdateUserFromID(db *gorm.DB, id string, data ProfileRequest) error {
	return db.Model(schema.User{}).Where("id = ?", id).Updates(data).Error
}

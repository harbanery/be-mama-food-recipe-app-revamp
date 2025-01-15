package account

import (
	"mama-recipe/schema"

	"gorm.io/gorm"
)

type AccountRepository interface {
	CheckUserID(db *gorm.DB, id string) error
	GetProfile(db *gorm.DB, id string) *Profile
	UpdateUser(db *gorm.DB, id string, data ProfileRequest) error
	UpdateUserPhoto(db *gorm.DB, id, url string) error
	DeleteUserPhoto(db *gorm.DB, id string) error
}

type accountRepository struct{}

func NewAccountRepository() AccountRepository {
	return &accountRepository{}
}

const idQueryParam = `id = ?`

func (s *accountRepository) CheckUserID(db *gorm.DB, id string) error {
	var user *schema.User
	return db.First(&user, idQueryParam, id).Error
}

func (s *accountRepository) GetProfile(db *gorm.DB, id string) *Profile {
	var user *Profile
	db.Model(&schema.User{}).First(&user, idQueryParam, id)
	return user
}

func (s *accountRepository) UpdateUser(db *gorm.DB, id string, data ProfileRequest) error {
	return db.Model(schema.User{}).Where(idQueryParam, id).Updates(data).Error
}

func (s *accountRepository) UpdateUserPhoto(db *gorm.DB, id, url string) error {
	return db.Model(schema.User{}).Where(idQueryParam, id).Update("photo", url).Error
}

func (s *accountRepository) DeleteUserPhoto(db *gorm.DB, id string) error {
	return db.Model(schema.User{}).Where(idQueryParam, id).Update("photo", nil).Error
}

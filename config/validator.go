package config

import (
	"mama-recipe/helper"

	"github.com/go-playground/validator/v10"
)

func NewValidator(env *helper.EnvLoad) *validator.Validate {
	return validator.New()
}

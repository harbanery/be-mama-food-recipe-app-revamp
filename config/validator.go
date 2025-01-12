package config

import "github.com/go-playground/validator/v10"

func NewValidator(env *EnvLoad) *validator.Validate {
	return validator.New()
}
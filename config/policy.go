package config

import (
	"mama-recipe/helper"

	"github.com/microcosm-cc/bluemonday"
)

func NewPolicy(env *helper.EnvLoad) *bluemonday.Policy {
	return bluemonday.UGCPolicy()
}

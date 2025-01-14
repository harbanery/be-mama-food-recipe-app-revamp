package config

import (
	"log"
	"mama-recipe/helper"

	"github.com/cloudinary/cloudinary-go/v2"
)

func NewCloudinary(env *helper.EnvLoad) *cloudinary.Cloudinary {
	cld, err := cloudinary.NewFromParams(env.CLOUDINARY_CLOUD_NAME, env.CLOUDINARY_API_KEY, env.CLOUDINARY_API_SECRET)
	if err != nil {
		log.Fatalf("failed to initialize cloudinary: %v", err)
	}
	return cld
}

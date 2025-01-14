package main

import (
	"log"
	"mama-recipe/config"
)

func main() {
	env := config.CheckEnv()

	db := config.NewDatabase(env)
	validate := config.NewValidator(env)
	policy := config.NewPolicy(env)
	app := config.NewFiber(env)
	cloudinary := config.NewCloudinary(env)

	config.Bootstrap(&config.BootstrapConfig{
		DB:          db,
		App:         app,
		Validate:    validate,
		Policy:      policy,
		Environment: env,
		Cloudinary:  cloudinary,
	})

	err := app.Listen(config.GetPort())
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

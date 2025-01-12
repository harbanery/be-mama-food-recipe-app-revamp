package main

import (
	"log"
	"mama-recipe/config"
)

func main() {
	env := config.CheckEnv()

	db := config.NewDatabase(env)
	validate := config.NewValidator(env)
	app := config.NewFiber(env)

	config.Bootstrap(&config.BootstrapConfig{DB: db, App: app, Validate: validate})

	err := app.Listen(config.GetPort())
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
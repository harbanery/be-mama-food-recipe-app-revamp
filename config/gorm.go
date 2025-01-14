package config

import (
	"fmt"
	"log"
	"mama-recipe/helper"
	"mama-recipe/schema"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(env *helper.EnvLoad) *gorm.DB {
	username := env.DB_USER
	password := env.DB_PASS
	host := env.DB_HOST
	port := env.DB_PORT
	database := env.DB_NAME
	require := env.DB_SSLMODE

	idleConnection, errIdleConnection := strconv.Atoi(env.DB_POOL_IDLE)
	maxConnection, errMaxConnection := strconv.Atoi(env.DB_POOL_MAX)
	maxLifeTimeConnection, errMaxLifeTimeConnection := strconv.Atoi(env.DB_POOL_LIFETIME)
	if errIdleConnection != nil || errMaxConnection != nil || errMaxLifeTimeConnection != nil {
		log.Fatalf("failed to convert string to integer")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta", host, username, password, database, port, require)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Silent,
				IgnoreRecordNotFoundError: true,
				ParameterizedQueries:      true,
				Colorful:                  false,
			},
		),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	connection, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get database connection: %v", err)
	}

	connection.SetMaxIdleConns(idleConnection)
	connection.SetMaxOpenConns(maxConnection)
	connection.SetConnMaxLifetime(time.Second * time.Duration(maxLifeTimeConnection))

	err = db.AutoMigrate(
		&schema.User{},
		&schema.Recipe{},
		&schema.Save{},
		&schema.Like{},
		// &schema.Comment{},
	)
	if err != nil {
		log.Fatalf("failed to auto migrate: %v", err)
	}

	return db
}

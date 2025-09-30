package db

import (
	"fmt"
	"log"

	"github.com/lekovv/go-crud-simple/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB(env *config.Env) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", env.DBHost, env.DBUser, env.DBPassword, env.DBName, env.DBPort, env.DBSSLMode, env.DBTimezone)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	DB.Logger = logger.Default.LogMode(logger.Info)

	fmt.Println("Connected to database")
}

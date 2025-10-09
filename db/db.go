package db

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/lekovv/go-web-mvp/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB *gorm.DB
}

func ConnectDB(env *config.Env) *Database {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		env.DBHost,
		env.DBUser,
		env.DBPassword,
		env.DBName,
		env.DBPort,
		env.DBSSLMode,
		env.DBTimezone)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	db.Logger = logger.Default.LogMode(logger.Info)

	fmt.Println("Connected to database")

	migrateURL := "file://migrations"
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s&timezone=%s",
		env.DBUser,
		env.DBPassword,
		env.DBHost,
		env.DBPort,
		env.DBName,
		env.DBSSLMode,
		env.DBTimezone,
	)

	m, err := migrate.New(migrateURL, dbURL)
	if err != nil {
		log.Fatal("Failed to initialize migrations", err.Error())
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Failed to apply migrations", err.Error())
	}
	fmt.Println("Applied migrations successfully")

	return &Database{DB: db}
}

package database

import (
	"fmt"
	"log"
	"revenue/config"
	"revenue/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	cfg := config.LoadConfig()
	dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBPort,
	)
	tempDB, err := gorm.Open(mysql.Open(dsnWithoutDB), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	tempDB.Exec("CREATE DATABASE IF NOT EXISTS " + cfg.DBName)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	runMigrations()
	fmt.Println("Database connected successfully!")
}

func runMigrations() {
	err := DB.AutoMigrate(
		&models.Customer{},
		&models.Product{},
		&models.Order{},
	)
	if err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}
}

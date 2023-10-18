package db

import (
	"derso.com/imersao-fullcycle/codepix-go/domain/model"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
)

func init() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	err := godotenv.Load(basepath + "/../../.env")

	if err != nil {
		log.Fatalf("Error loading .env files")
	}
}

func ConnectDB(env string) *gorm.DB {
	var config *gorm.Config
	var dsn string
	var db *gorm.DB
	var err error

	if os.Getenv("debug") == "true" {
		config = &gorm.Config{Logger: logger.Default.LogMode(logger.Info)}
	} else {
		config = &gorm.Config{Logger: logger.Default.LogMode(logger.Error)}
	}

	if env != "test" {
		dsn = os.Getenv("dsn")
		db, err = gorm.Open(postgres.Open(dsn), config)
	} else {
		dsn = os.Getenv("dsnTest")
		db, err = gorm.Open(sqlite.Open(dsn), config)
	}

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
		panic(err)
	}

	if os.Getenv("AutoMigrateDb") == "true" {
		db.AutoMigrate(&model.Bank{}, &model.Account{}, &model.PixKey{}, &model.Transaction{})
	}

	return db
}
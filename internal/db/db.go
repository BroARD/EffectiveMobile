package db

import (
	personservice "EffectiveMobile/internal/personService"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func init(){
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func buildDSN() string {
    return "host=" + os.Getenv("DB_HOST") +
        " port=" + os.Getenv("DB_PORT") +
        " user=" + os.Getenv("DB_USER") +
        " password=" + os.Getenv("DB_PASSWORD") +
        " dbname=" + os.Getenv("DB_NAME") +
        " sslmode=" + os.Getenv("DB_SSLMODE")
}

func InitDB() (*gorm.DB, error){
	dsn := buildDSN()
	var err error

	log.Println("INFO: Инициализация подключения к базе данных")
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("FATAL: Ошибка подключения к БД: %v", err)
	}

	log.Println("INFO: Применение миграций")
	if err := db.AutoMigrate(&personservice.Person{}); err != nil {
		log.Fatalf("FATAL: Ошибка миграции: %v", err)
	}
	log.Println("INFO: База данных готова")

	return db, nil
}

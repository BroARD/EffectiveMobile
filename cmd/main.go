package main

import (
	"EffectiveMobile/internal/db"
	"EffectiveMobile/internal/handlers"
	"EffectiveMobile/internal/personService"
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/webstradev/echo-pagination/pkg/pagination"

	"github.com/swaggo/echo-swagger"
	_ "EffectiveMobile/docs"
)

// @title Person App API
// @version 1.0

// @host localhost:8080
// @BasePath /

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	fmt.Println("INFO: Запуск сервера")
	database , err := db.InitDB()
	if err != nil {
		log.Fatal("Could not connect to DB")
	}

	persRepo := personservice.NewPersonRepository(database)
	persService := personservice.NewPersonService(persRepo)
	persHandlers := handlers.NewPersonHandler(persService)

	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "INFO: method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.Use(pagination.New(pagination.WithMinPageSize(1)))
	e.Validator = &CustomValidator{validator: validator.New()}

	e.GET("/persons", persHandlers.GetPersons)
	e.POST("/person", persHandlers.PostPerson)
	e.DELETE("/person/:person_id", persHandlers.DelPerson)
	e.PATCH("/person/:person_id", persHandlers.UpdatePerson)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	log.Println("INFO: Сервер запущен на localhost:8080")
	e.Start("localhost:8080")
}
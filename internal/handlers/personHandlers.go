package handlers

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"EffectiveMobile/internal/personService"
)

type PersHandler struct {
	service personservice.PersonService
}

func NewPersonHandler(service personservice.PersonService) *PersHandler {
	return &PersHandler{service: service}
}

// @Summary GetPersons
// @Tags Get
// @Descriptions get all users
// @ID get-persons
// @Param page query int false "Номер страницы" default(1)
// @Param size query int false "Количество записей" default(10)
// @Param age query int false "Возраст"
// @Param gender query string false "Пол"
// @Param country query string false "Национальность"
// @Accept json
// @Produce json
// @Router /persons [get]
func (h *PersHandler) GetPersons(c echo.Context) error {
	log.Println("DEBUG: Начало обработки запроса GET /persons")

	pageParam := c.QueryParam("page")
	sizeParam := c.QueryParam("size")
	ageParam := c.QueryParam("age")
	genderParam := c.QueryParam("gender")
	countryParam := c.QueryParam("country")

	persons, err := h.service.GetAllPersons(pageParam, sizeParam, ageParam, genderParam, countryParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, persons)
}

// @Summary DeletePerson
// @Tags Delete
// @Descriptions Delete user by ID
// @ID delete-person
// @Accept json
// @Produce json
// @Param person_id path string true "person_id"
// @Router /person/{person_id} [delete]
func (h *PersHandler) DelPerson(c echo.Context) error {
	personID := c.Param("person_id")

	if _, err := h.service.GetPersonByID(personID); err != nil{
		log.Printf("ERROR: Пользователь не найден: %v", err)
		return c.JSON(http.StatusNotFound, map[string]string{personID:"Пользователь с данным не найден"})
	}

	log.Printf("INFO: Попытка удаления пользователя ID: %s", personID)
	if err := h.service.DeletePerson(personID); err != nil {
		log.Printf("ERROR: Ошибка удаления: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	log.Printf("INFO: Успешно удален пользователь ID: %s", personID)
	return c.JSON(http.StatusOK, map[string]string{personID:"Пользователь с данным ID удалён"})
}

// @Summary UpdatePerson
// @Tags Update
// @Descriptions Update user by ID
// @ID update-person
// @Accept json
// @Produce json
// @Param person_id path string true "person_id"
// @Param input body personservice.Person true "person info"
// @Router /person/{person_id} [patch]
func (h *PersHandler) UpdatePerson(c echo.Context) error {
	personID := c.Param("person_id")
	log.Printf("INFO: Начало обновления пользователя ID: %s", personID)

	var personNew personservice.Person

	if err := c.Bind(&personNew); err != nil {
		log.Printf("ERROR: Ошибка парсинга JSON: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Неверный формат данных"})
	}
	log.Printf("DEBUG: Получены данные для обновления: %+v", personNew)

	personUpdate, err := h.service.UpdatePeson(personID, personNew)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, personUpdate)
}

// @Summary CreatePerson
// @Tags Create
// @Descriptions Create person
// @ID create-person
// @Accept json
// @Produce json
// @Param input body personservice.Person true "person info"
// @Router /person [post]
func (h *PersHandler) PostPerson(c echo.Context) error {
	log.Println("INFO: Начало создания нового пользователя")
	
	var person personservice.Person
	if err := c.Bind(&person); err != nil {
		log.Printf("ERROR: Ошибка парсинга JSON: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Неверный формат данных"})
	}

	if err := c.Validate(&person); err != nil {
		log.Printf("WARN: Ошибка валидации: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Не заполнены обязательные поля"})
	}
	log.Println("DEBUG: Валидация пройдена успешно")

	log.Printf("DEBUG: Получение дополнительных данных для: %s %s", person.Name, person.Surname)
	
	newPerson, err := h.service.CreatePerson(person)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	log.Printf("INFO: Создан новый пользователь ID: %s", person.ID)
	return c.JSON(http.StatusCreated, newPerson)
}
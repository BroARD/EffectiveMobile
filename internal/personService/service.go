package personservice

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/uuid"
	"github.com/gouef/country"
)

type PersonService interface {
	CreatePerson(person Person) (Person, error)
	GetAllPersons(pageParam string, sizeParam string, ageParam string, genderParam string, countryParam string) ([]Person, error)
	GetPersonByID(personID string) (Person, error)
	UpdatePeson(personID string, person Person) (Person, error)
	DeletePerson(personID string) error
}

type persService struct {
	repo PersonRepository
}

func NewPersonService(r PersonRepository) PersonService {
	return &persService{repo: r}
}

// CreatePerson implements PersonService.
func (s *persService) CreatePerson(person Person) (Person,  error) {
	age, err := s.GetAgeByName(person.Name)
	if err != nil {
		return Person{}, err
	}
	gender, err := s.GetGenderByName(person.Name)
	if err != nil {
		return Person{}, err
	}
	country, err := s.GetCountryByName(person.Name)
	if err != nil {
		return Person{}, err
	}

	person.ID = uuid.NewString()
	person.Age = age
	person.Gender = gender
	person.Country = country

	if err := s.repo.CreatePerson(person); err != nil {
		return Person{}, err
	}

	return person, nil
}

// GetAllPersons implements PersonService.
func (s *persService) GetAllPersons(pageParam string, sizeParam string, ageParam string, genderParam string, countryParam string) ([]Person, error) {
	page := 1
	size := 10
	a := 0

	if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
		page = p
	}
	if s, err := strconv.Atoi(sizeParam); err == nil && s > 0 {
		size = s
	}

	offset := (page - 1) * size

	if ageParam != "" {
		a, _ = strconv.Atoi(ageParam)
	} 
	
	return s.repo.GetAllPersons(size, a, genderParam, countryParam, offset)
}

func (s *persService) DeletePerson(personID string) error {
	return s.repo.DeletePerson(personID)
}

func (s *persService) GetPersonByID(personID string) (Person, error) {
	return s.repo.GetPersonByID(personID)
}

func (s *persService) UpdatePeson(personID string, person Person) (Person, error) {
	person_old, err := s.repo.GetPersonByID(personID)
	if err != nil {
		return Person{}, err
	}

	if err := s.repo.UpdatePerson(person_old, person); err != nil {
		return Person{}, err
	}

	return person, nil
}

func (s *persService) GetAgeByName(name string) (int, error) {
	log.Printf("DEBUG: Запрос возраста для имени: %s", name)

	baseURL := "https://api.agify.io/"
	params := url.Values{}
	params.Add("name", name)

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		log.Printf("ERROR: Ошибка HTTP-запроса к agify: %v", err)
		return 0, err
	}
	defer resp.Body.Close()

	bs := make([]byte, 1024)
	n, _ := resp.Body.Read(bs)

	var result StructForGetAge
	if err := json.Unmarshal(bs[:n], &result); err != nil {
		log.Printf("ERROR: Ошибка парсинга ответа agify: %v", err)
		return 0, err
	}

	log.Printf("DEBUG: Получен возраст %d для имени %s", result.Age, name)
	return result.Age, nil
}

func (s *persService) GetGenderByName(name string) (string, error) {
	log.Printf("DEBUG: Запрос пола для имени: %s", name)

	baseURL := "https://api.genderize.io/"
	params := url.Values{}
	params.Add("name", name)

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		log.Printf("ERROR: Ошибка HTTP-запроса к genderize: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	bs := make([]byte, 1024)
	n, _ := resp.Body.Read(bs)

	var result StructForGetGender
	if err := json.Unmarshal(bs[:n], &result); err != nil {
		log.Printf("ERROR: Ошибка парсинга ответа genderize: %v", err)
		return "", err
	}

	log.Printf("DEBUG: Получен пол %s для имени %s", result.Gender, name)
	return result.Gender, nil
}

func (s *persService) GetCountryByName(name string) (string, error) {
	log.Printf("DEBUG: Запрос страны для имени: %s", name)

	baseURL := "https://api.nationalize.io/"
	params := url.Values{}
	params.Add("name", name)

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		log.Printf("ERROR: Ошибка HTTP-запроса к nationalize: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	bs := make([]byte, 1024)
	n, _ := resp.Body.Read(bs)

	var result StructForGetCountry
	if err := json.Unmarshal(bs[:n], &result); err != nil {
		log.Printf("ERROR: Ошибка парсинга ответа nationalize: %v", err)
		return "", err
	}

	if len(result.Country) == 0 {
		log.Printf("WARN: Нет данных о стране для имени %s", name)
		return "Unknown", nil
	}

	countryCode := result.Country[0]["country_id"].(string)
	countryInfo := country.FindByAlpha2(countryCode)
	if countryInfo == nil {
		log.Printf("WARN: Неизвестный код страны: %s", countryCode)
		return "Unknown", nil
	}

	log.Printf("DEBUG: Получена страна %s для имени %s", countryInfo.Name, name)
	return countryInfo.Name, nil
}
package personservice

import "gorm.io/gorm"

// Основные методы CRUD

type PersonRepository interface {
	CreatePerson(person Person) error
	GetAllPersons(sizeParam int, ageParam int, genderParam string, countryParam string, offset int) ([]Person, error)
	GetPersonByID(personID string) (Person, error)
	UpdatePerson(person_old Person, person_new Person) error
	DeletePerson(personID string) error
}

type persRepository struct {
	db *gorm.DB
}

func NewPersonRepository(db *gorm.DB) PersonRepository {
	return &persRepository{db: db}
}

func (r *persRepository) CreatePerson(person Person) error {
	return r.db.Create(&person).Error
}

func (r *persRepository) GetAllPersons(sizeParam int, ageParam int, genderParam string, countryParam string, offset int) ([]Person, error) {
	var persons []Person

	err := r.db.Where(&Person{Gender: genderParam, Age: ageParam, Country: countryParam}).Offset(offset).Limit(sizeParam).Find(&persons).Error
	return persons, err
}

func (r *persRepository) GetPersonByID(personID string) (Person, error) {
	var person Person

	err := r.db.First(&person, "id = ?", personID).Error
	return person, err
}

func (r *persRepository) UpdatePerson(person_old Person, person_new Person) error {
	return r.db.Model(&person_old).Updates(&person_new).Error
}

func (r *persRepository) DeletePerson(personID string) error{
	return r.db.Where("id = ?", personID).Delete(&Person{}).Error
}

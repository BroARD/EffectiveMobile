package personservice

type Person struct {
	ID         string `json:"id" gorm:"primaryKey"`
	Name       string `json:"name" validate:"required"`
	Surname    string `json:"surname" validate:"required"`
	Patronymic string `json:"patronymic"`
	Age        int    `json:"age"`
	Gender     string `json:"gender"`
	Country    string `json:"country"`
}

type (
	StructForGetAge struct {
		Count int    `json:"count"`
		Name  string `json:"name"`
		Age   int    `json:"age"`
	}

	StructForGetCountry struct {
		Count   int              `json:"count"`
		Name    string           `json:"name"`
		Country []map[string]any `json:"country"`
	}

	StructForGetGender struct {
		Count      int     `json:"count"`
		Name       string  `json:"name"`
		Gender     string  `json:"gender"`
		Probability float64 `json:"probability"`
	}
)
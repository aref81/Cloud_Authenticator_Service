package model

type User struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	NationalCode string `json:"national_code"`
	IPAddress    string `json:"ip"`
	Status       string `json:"status"`
}

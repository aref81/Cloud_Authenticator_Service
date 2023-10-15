package model

type User struct {
	Email        string `json:"email"`
	LastName     string `json:"last_name"`
	NationalCode string `json:"national_code"`
	IPAddress    string `json:"-"`
	Photo1       string `json:"photo1"`
	Photo2       string `json:"photo2"`
	Username     string `json:"-"`
}

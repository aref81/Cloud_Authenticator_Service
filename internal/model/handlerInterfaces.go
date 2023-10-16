package model

// Register Handler
type RegisterReq struct {
	Name         string `json:"name" validate:"required,max=40"`
	Email        string `json:"email" validate:"required,email"`
	NationalCode string `json:"national_code" validate:"required,numeric,max=10"`
	Pic1         string `json:"pic_1" validate:"required"`
	Pic2         string `json:"pic_2" validate:"required"`
	IPAddress    string `json:"ip"`
}

type RegisterRes struct {
	Message string `json:"message"`
}

// Status Handler
type StatusRes struct {
	Status string `json:"status"`
}

type StatusReq struct {
	NationalCode string `json:"national_code" validate:"required,numeric,max=10"`
}

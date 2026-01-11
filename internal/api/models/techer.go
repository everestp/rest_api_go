package models

type Teacher struct{
	ID int `json:"id"`
	FirstNAme string  `json:"first_name"`
	LastName string  `json:"last_name"`
	Class string `json:"class"`
	Subject string `json:"subject"`

}
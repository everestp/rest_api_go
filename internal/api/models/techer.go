package models

type Teacher struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Level     string `json:"level"`
	Subject   string `json:"subject"`
}

package models

type User struct {
	Model
	Name           string `json:"name"`
	PhoneNumber    string `json:"phone_number"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	HashedPassword string `json:"hashed_password"`
	IsAgree        bool   `json:"is_agree"`
}

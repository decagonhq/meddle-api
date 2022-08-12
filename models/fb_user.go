package models

const alphabet = "abcdefghhijklmnopqrstuvwxyz"

type FacebookUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

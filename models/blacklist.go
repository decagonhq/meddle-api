package models

import "time"

type BlackList struct {
	Model
	Token     string `json:"token"`
	Email     string `json:"email"`
	CreatedAt time.Time
}

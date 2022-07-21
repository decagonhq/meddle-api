package dto

type SignupResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

type SignupRequest struct {
	Name        string `json:"name" gorm:"not null" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required,e164" gorm:"not null" gorm:"unique"`
	Email       string `json:"email" binding:"required,email" gorm:"not null" gorm:"unique"`
	Password    string `json:"password" binding:"required" gorm:"not null"`
}

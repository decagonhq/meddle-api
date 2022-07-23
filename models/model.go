package models

type Model struct {
	ID        string `json:"id" gorm:"primaryKey"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	DeletedAt int64  `json:"deleted_at"`
}

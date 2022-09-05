package models

type FCMNotificationToken struct {
	Model
	UserID   uint   `gorm:"column:user_id"`
	Token    string `gorm:"column:token"`
	isViewed bool   `gorm:"column:is_viewed"`
}

type AddNotificationTokenArgs struct {
	Token  string `json:"token" binding:"required"`
	UserID uint   `json:"user_id"`
}

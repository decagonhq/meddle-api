package models

type PushNotificationCategory string

const (
	NextMedicationCategory PushNotificationCategory = "NEXT_MEDICATION_CATEGORY"
	WelcomeCategory        PushNotificationCategory = "WELCOME_CATEGORY"
)

type FCMNotificationToken struct {
	Model
	UserID   uint   `json:"user_id"`
	Token    string `json:"token"`
	IsViewed bool   `json:"is_viewed"`
}

type AddNotificationTokenArgs struct {
	Token  string `json:"token" binding:"required"`
	UserID uint   `json:"user_id"`
}

type PushPayload struct {
	Title       string                   `json:"title"`
	Body        string                   `json:"body"`
	Data        map[string]string        `json:"data"`
	ClickAction string                   `json:"clickAction"`
	Category    PushNotificationCategory `json:"category"`
}

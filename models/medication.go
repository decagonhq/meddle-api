package models

import (
	"gorm.io/gorm"
	"time"
)

type Medication struct {
	//base model goes here
	gorm.Model
	Name                   string    `gorm:"column:name"`
	Dosage                 int       `gorm:"column:dosage"`
	TimeInterval           int       `gorm:"column:time_interval"` // min hour daily
	MedicationStartDate    time.Time `gorm:"column:medication_start_date"`
	Duration               int       `gorm:"column:duration"`
	MedicationPrescribedBy string    `gorm:"column:medication_prescribed_by"`
	MedicationStopDate     time.Time `gorm:"column:medication_stop_date"`
	MedicationStartTime    time.Time `gorm:"column:medication_start_time"`
	NextDosageTime         time.Time `gorm:"column:next_dosage_time"`
	PurposeOfMedication    string    `gorm:"column:purpose_of_medication"`
	IsMedicationDone       bool      `gorm:"column:is_medication_done"`
	UserID                 uint      `gorm:"column:user_id"`
}

type MedicationRequest struct {
	Name                   string `json:"name" binding:"required"`
	Dosage                 int    `json:"dosage" binding:"required"`
	TimeInterval           int    `json:"time_interval" binding:"required"` // min hour daily
	MedicationStartDate    string `json:"medication_start_date" binding:"required"`
	Duration               int    `json:"duration" binding:"required"`
	MedicationPrescribedBy string `json:"medication_prescribed_by" binding:"required"`
	MedicationStopDate     string `json:"medication_stop_date" binding:"required"`
	MedicationStartTime    string `json:"medication_start_time" binding:"required"`
	PurposeOfMedication    string `json:"purpose_of_medication" binding:"required"`
	UserID                 uint   `json:"user_id"`
}

type MedicationResponse struct {
	ID                     uint   `json:"id"`
	CreatedAt              string `json:"created_at"`
	UpdatedAt              string `json:"updated_at"`
	Name                   string `json:"name"`
	Dosage                 int    `json:"dosage"`
	TimeInterval           int    `json:"time_interval"` // min hour daily
	MedicationStartDate    string `json:"medication_start_date"`
	Duration               int    `json:"duration"`
	MedicationPrescribedBy string `json:"medication_prescribed_by"`
	MedicationStopDate     string `json:"medication_stop_date"`
	MedicationStartTime    string `json:"medication_start_time"`
	NextDosageTime         string `json:"next_dosage_time"`
	PurposeOfMedication    string `json:"purpose_of_medication"`
	UserID                 uint   `json:"user_id"`
}

func (m *MedicationRequest) ReqToMedicationModel() *Medication {
	return &Medication{
		Name:                   m.Name,
		Dosage:                 m.Dosage,
		TimeInterval:           m.TimeInterval,
		Duration:               m.Duration,
		MedicationPrescribedBy: m.MedicationPrescribedBy,
		PurposeOfMedication:    m.PurposeOfMedication,
		UserID:                 m.UserID,
	}
}

func (m *Medication) MedicationToResponse() *MedicationResponse {
	return &MedicationResponse{
		ID:                     m.ID,
		CreatedAt:              m.CreatedAt.String(),
		UpdatedAt:              m.UpdatedAt.String(),
		Name:                   m.Name,
		Dosage:                 m.Dosage,
		TimeInterval:           m.TimeInterval,
		MedicationStartDate:    m.MedicationStartDate.String(),
		Duration:               m.Duration,
		MedicationPrescribedBy: m.MedicationPrescribedBy,
		MedicationStopDate:     m.MedicationStopDate.String(),
		MedicationStartTime:    m.MedicationStartTime.String(),
		NextDosageTime:         m.NextDosageTime.String(),
		PurposeOfMedication:    m.PurposeOfMedication,
		UserID:                 m.UserID,
	}
}

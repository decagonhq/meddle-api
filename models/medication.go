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
	gorm.Model
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
	gorm.Model
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
	startDate := m.MedicationStartDate.String()
	stopDate := m.MedicationStopDate.String()
	startTime := m.MedicationStartTime.String()
	return &MedicationResponse{
		Model:                  gorm.Model{},
		Name:                   m.Name,
		Dosage:                 m.Dosage,
		TimeInterval:           m.TimeInterval,
		MedicationStartDate:    startDate,
		Duration:               m.Duration,
		MedicationPrescribedBy: m.MedicationPrescribedBy,
		MedicationStopDate:     stopDate,
		MedicationStartTime:    startTime,
		PurposeOfMedication:    m.PurposeOfMedication,
		UserID:                 m.UserID,
	}
}

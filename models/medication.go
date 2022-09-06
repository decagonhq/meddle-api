package models

import (
	"time"
)

type Medication struct {
	//base model goes here
	Model
	Name                   string    `json:"name"`
	Dosage                 int       `json:"dosage"`
	TimeInterval           int       `json:"time_interval"` // min hour daily
	MedicationStartDate    time.Time `json:"medication_start_date"`
	Duration               int       `json:"duration"`
	MedicationPrescribedBy string    `json:"medication_prescribed_by"`
	MedicationStopDate     time.Time `json:"medication_stop_date"`
	MedicationStartTime    time.Time `json:"medication_start_time"`
	NextDosageTime         time.Time `json:"next_dosage_time"`
	PurposeOfMedication    string    `json:"purpose_of_medication"`
	IsMedicationDone       bool      `json:"is_medication_done"`
	MedicationIcon         string    `json:"medication_icon"`
	UserID                 uint      `json:"user_id"`
}

type UpdateMedicationRequest struct {
	Name                   string `json:"name"`
	Dosage                 int    `json:"dosage"`
	TimeInterval           int    `json:"time_interval"` // min hour daily
	MedicationStartDate    string `json:"medication_start_date"`
	Duration               int    `json:"duration"`
	MedicationPrescribedBy string `json:"medication_prescribed_by"`
	MedicationStopDate     string `json:"medication_stop_date"`
	MedicationStartTime    string `json:"medication_start_time"`
	PurposeOfMedication    string `json:"purpose_of_medication"`
	MedicationIcon         string `json:"medication_icon"`
}

type MedicationRequest struct {
	Name                   string `json:"name" binding:"required"`
	Dosage                 int    `json:"dosage" binding:"required"`
	TimeInterval           int    `json:"time_interval" binding:"required"` // min hour daily
	MedicationStartDate    string `json:"medication_start_date" binding:"required"`
	Duration               int    `json:"duration" binding:"required"`
	MedicationPrescribedBy string `json:"medication_prescribed_by" binding:"required"`
	MedicationStartTime    string `json:"medication_start_time" binding:"required"`
	PurposeOfMedication    string `json:"purpose_of_medication" binding:"required"`
	MedicationIcon         string `json:"medication_icon" binding:"required"`
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
	MedicationIcon         string `json:"medication_icon"`
	UserID                 uint   `json:"user_id"`
}

type MedicationDetailResponse struct {
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
	MedicationIcon         string `json:"medication_icon"`
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
		MedicationIcon:         m.MedicationIcon,
		UserID:                 m.UserID,
	}
}

func (m *Medication) MedicationToResponse() *MedicationResponse {
	return &MedicationResponse{
		ID:                     m.ID,
		CreatedAt:              time.Unix(m.CreatedAt, 0).String(),
		UpdatedAt:              time.Unix(m.UpdatedAt, 0).String(),
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
		MedicationIcon:         m.MedicationIcon,
		UserID:                 m.UserID,
	}
}

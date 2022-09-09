package models

import "time"

type MedicationHistory struct {
	Model
	MedicationName         string    `json:"medication_name"`
	MedicationID           uint      `json:"medication_id"`
	MedicationTime         time.Time `json:"medication_time"`
	MedicationDosage       int       `json:"medication_dosage"`
	UserID                 uint      `json:"user_id"`
	HasMedicationBeenTaken bool      `json:"has_medication_been_taken"`
}

func NewMedicationHistory(medication Medication) *MedicationHistory {
	return &MedicationHistory{
		MedicationName:         medication.Name,
		MedicationID:           medication.ID,
		MedicationDosage:       medication.Dosage,
		MedicationTime:         medication.NextDosageTime,
		UserID:                 medication.UserID,
		HasMedicationBeenTaken: false,
	}

}

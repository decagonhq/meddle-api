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
	WasMedicationMissed    string    `json:"was_medication_missed"`
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

type MedicationHistoryResponse struct {
	ID                     uint   `json:"id"`
	CreatedAt              string `json:"created_at"`
	UpdatedAt              string `json:"updated_at"`
	MedicationName         string `json:"medication_name"`
	MedicationID           uint   `json:"medication_id"`
	MedicationTime         string `json:"medication_time"`
	MedicationDosage       int    `json:"medication_dosage"`
	UserID                 uint   `json:"user_id"`
	HasMedicationBeenTaken bool   `json:"has_medication_been_taken"`
	WasMedicationMissed    string `json:"was_medication_missed"`
}

func (m *MedicationHistory) MedicationHistoryToResponse() *MedicationHistoryResponse {
	return &MedicationHistoryResponse{
		ID:                     m.ID,
		CreatedAt:              time.Unix(m.CreatedAt, 0).String(),
		UpdatedAt:              time.Unix(m.UpdatedAt, 0).String(),
		MedicationName:         m.MedicationName,
		MedicationID:           m.MedicationID,
		MedicationTime:         m.MedicationTime.UTC().String(),
		MedicationDosage:       m.MedicationDosage,
		UserID:                 m.UserID,
		HasMedicationBeenTaken: m.HasMedicationBeenTaken,
		WasMedicationMissed:    m.WasMedicationMissed,
	}
}

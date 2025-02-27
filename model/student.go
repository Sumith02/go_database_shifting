package model

type StudentData struct {
	StudentUnitID   string `json:"student_unit_id"`
	StudentName     string `json:"student_name"`
	StudentUSN      string `json:"student_usn"`
	Department      string `json:"department"`
	UnitID          string `json:"unit_id"`
	FingerprintData string `json:"fingerprint_data"`
}

package repository

import (
	"database/sql"
	"db_shifter/model"
	"log"
)

type StudentRepository struct {
	DB *sql.DB
}

func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{DB: db}
}

func (r *StudentRepository) FetchStudents(machineId string) ([]model.StudentData, error) {
	query := `SELECT v.student_unit_id, v.student_name, v.student_usn, v.department, f.unit_id, f.fingerprint
			  FROM ` + machineId + ` v
			  JOIN fingerprintdata f ON v.student_id = f.student_id;`

	rows, err := r.DB.Query(query)
	if err != nil {
		log.Println("Error fetching data:", err)
		return nil, err
	}
	defer rows.Close()

	var students []model.StudentData

	for rows.Next() {
		var student model.StudentData
		err := rows.Scan(&student.StudentUnitID, &student.StudentName, &student.StudentUSN, &student.Department, &student.UnitID, &student.FingerprintData)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	return students, nil
}

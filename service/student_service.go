package service

import (
	service "db_shifter/api"
	"db_shifter/repository"
	"log"
)

type StudentService struct {
	Repo *repository.StudentRepository
}

func NewStudentService(repo *repository.StudentRepository) *StudentService {
	return &StudentService{Repo: repo}
}

func (s *StudentService) ProcessAndTransferStudents(machineId string, apiClient *service.APIClient) error {
	students, err := s.Repo.FetchStudents(machineId)
	if err != nil {
		return err
	}

	for _, student := range students {
		err := apiClient.SendStudentData(student)
		if err != nil {
			log.Printf("Error transferring student %s: %v", student.StudentName, err)
		} else {
			log.Printf("Successfully transferred student: %s", student.StudentName)
		}
	}

	return nil
}

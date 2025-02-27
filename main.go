package main

import (
	"db_shifter/api"
	"db_shifter/config"
	"db_shifter/repository"
	"db_shifter/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	// Connect to Old Database
	db, err := config.ConnectOldDB()
	if err != nil {
		log.Fatal("Database connection error:", err)
	}
	defer db.Close()

	// Initialize layers
	studentRepo := repository.NewStudentRepository(db)
	studentService := service.NewStudentService(studentRepo)
	apiClient := api.NewAPIClient()

	// Process and Transfer Data
	router := mux.NewRouter()

	router.HandleFunc("/{machineId}", func(w http.ResponseWriter, r *http.Request) {
		machineId := mux.Vars(r)["machineId"]
		err = studentService.ProcessAndTransferStudents(machineId, apiClient)
		if err != nil {
			log.Fatal("Error in processing:", err)
		}
	})

	log.Println("server is running..")

	http.ListenAndServe(":8080", router)
}

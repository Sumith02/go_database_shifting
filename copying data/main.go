package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

type StudentPayload struct {
	StudentUnitId   string `json:"student_unit_id"`
	StudentName     string `json:"student_name"`
	StudentUsn      string `json:"student_usn"`
	Department      string `json:"department"`
	UnitId          string `json:"unit_id"`
	FingerprintData string `json:"fingerprint_data"`
}

func main() {
	if len(os.Args) < 3 {
		log.Fatal("❌ Usage: go run migrate_students.go <source_table_name> <target_unit_id>")
	}

	sourceTable := os.Args[1] // e.g., vs242s25
	targetUnitId := os.Args[2] // e.g., vs242s101

	// PostgreSQL connection
	connStr := "postgres://postgres:B8evHxztKxT6nrqI8OgF@vithsutra-production.chkeii8oa8ak.eu-north-1.rds.amazonaws.com:5432/vithsutra_biometric_db1"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("❌ DB connection failed: %v", err)
	}
	defer db.Close()

	// SQL query to extract from given source table and join with fingerprintdata
	query := fmt.Sprintf(`
		SELECT 
			s.student_unit_id,
			s.student_name,
			s.student_usn,
			s.department,
			f.fingerprint
		FROM %s s
		JOIN fingerprintdata f 
		  ON s.student_id = f.student_id AND s.student_unit_id = f.student_unit_id
	`, pqQuoteIdentifier(sourceTable))

	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("❌ Query failed: %v", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var student StudentPayload

		err := rows.Scan(
			&student.StudentUnitId,
			&student.StudentName,
			&student.StudentUsn,
			&student.Department,
			&student.FingerprintData,
		)
		if err != nil {
			log.Printf("⚠️ Row scan error: %v\n", err)
			continue
		}

		student.UnitId = targetUnitId // Inject unit_id dynamically

		body, _ := json.Marshal(student)

		resp, err := http.Post(
			"https://biometric.http.vsensetech.in/users/newstudent",
			"application/json",
			bytes.NewBuffer(body),
		)

		if err != nil {
			log.Printf("❌ API call failed for %s: %v\n", student.StudentUnitId, err)
			continue
		}

		var result map[string]string
		json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			log.Printf("✅ Inserted: %s", student.StudentUnitId)
			count++
		} else {
			log.Printf("⚠️ Skipped %s: %s", student.StudentUnitId, result["error"])
		}
	}

	log.Printf("🎉 Done. Total inserted: %d", count)
}

// pqQuoteIdentifier safely quotes table name to prevent SQL injection
func pqQuoteIdentifier(id string) string {
	return `"` + strings.ReplaceAll(id, `"`, `""`) + `"`
}

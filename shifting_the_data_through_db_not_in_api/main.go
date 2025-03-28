package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// Define connection strings
	const oldDBConnStr = "postgresql://postgres:B8evHxztKxT6nrqI8OgF@vithsutra-production.chkeii8oa8ak.eu-north-1.rds.amazonaws.com:5432/vithsutra_biometric_db1"
	const newDBConnStr = "postgresql://postgres:B8evHxztKxT6nrqI8OgF@vithsutra-production.chkeii8oa8ak.eu-north-1.rds.amazonaws.com:5432/vithsutra_biometric_db1"

	// Connect to the source database
	sourceDB, err := sql.Open("postgres", oldDBConnStr)
	if err != nil {
		log.Fatal("Error connecting to source database:", err)
	}
	defer sourceDB.Close()

	// Connect to the destination database
	destDB, err := sql.Open("postgres", newDBConnStr)
	if err != nil {
		log.Fatal("Error connecting to destination database:", err)
	}
	defer destDB.Close()

	// Query to fetch data from source database
	rows, err := sourceDB.Query("SELECT student_unit_id, unit_id, fingerprint FROM fingerprintdata WHERE unit_id=$1", "vs242s02")
	if err != nil {
		log.Fatal("Error fetching data from source database:", err)
	}
	defer rows.Close()

	// Prepare statement for inserting into destination database
	stmt, err := destDB.Prepare("INSERT INTO inserts (unit_id, student_unit_id, fingerprint_data) VALUES ($1, $2, $3)")
	if err != nil {
		log.Fatal("Error preparing insert statement:", err)
	}
	defer stmt.Close()

	// Iterate over the fetched rows and insert into the destination database
	for rows.Next() {
		var studentUnitID, unitID, fingerprint string
		if err := rows.Scan(&studentUnitID, &unitID, &fingerprint); err != nil {
			log.Fatal("Error scanning row:", err)
		}

		_, err = stmt.Exec(unitID, studentUnitID, fingerprint)
		if err != nil {
			log.Fatal("Error inserting data into destination database:", err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal("Error during iteration:", err)
	}

	fmt.Println("Data successfully copied!")
}

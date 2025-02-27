package config

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

const oldDBConnStr = "postgresql://postgres:B9PD0Ugyd8t9NCYEmmIX@vithsutra-db.chkeii8oa8ak.eu-north-1.rds.amazonaws.com:5432/biometric"

func ConnectOldDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", oldDBConnStr)
	if err != nil {
		log.Fatal("Failed to connect to Old Database:", err)
		return nil, err
	}
	return db, nil
}

package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Seat struct {
	ID             int
	SeatIdentifier string
	SeatClass      string
}

func main() {
	// Database connection string
	dsn := "host=your_host port=your_port user=your_username password=your_password dbname=your_database sslmode=require"

	// Open the database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Auto-migrate the Seat model to create the table if it doesn't exist
	err = db.AutoMigrate(&Seat{})
	if err != nil {
		log.Fatal(err)
	}

	// Open the CSV file
	file, err := os.Open("seats.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(bufio.NewReader(file))

	// Read the CSV records
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Iterate over the records and insert them into the database
	for _, record := range records {
		// Parse the record values
		id, err := strconv.Atoi(record[0])
		if err != nil {
			log.Printf("Invalid ID for record: %v", record)
			continue
		}

		seat := Seat{
			ID:             id,
			SeatIdentifier: record[1],
			SeatClass:      record[2],
		}

		// Insert the seat into the database
		result := db.Create(&seat)
		if result.Error != nil {
			log.Printf("Error inserting seat: %v", result.Error)
			continue
		}

		log.Printf("Seat inserted successfully: ID=%d", seat.ID)
	}

	fmt.Println("CSV data inserted into the database.")
}

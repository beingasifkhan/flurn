package database

import (
	"bufio"
	"encoding/csv"
	"flurn_assignment/models"
	"fmt"
	"log"
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	dsn := "host=localhost user=postgres password=25205089p dbname=bookings_app port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	db.AutoMigrate(&models.Seat{}, &models.SeatPricing{}, &models.Booking{})
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

		seat := models.Seat{
			ID:             uint(id),
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

	DB = db
	return nil
}

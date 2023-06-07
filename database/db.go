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
	dropTables(db)

	db.AutoMigrate(&models.Seat{}, &models.SeatPricing{}, &models.Users{}, &models.Bookings{})
	if err != nil {
		log.Fatal(err)
	}

	users := []models.Users{
		{UserName: "Asif", Email: "beingasifkhan17@gmail.com", PhoneNumber: "7977148296"},
		{UserName: "Salman", Email: "asnx@gmail.com", PhoneNumber: "9619866554"},
		{UserName: "Shah Rukh", Email: "beingsana@gamil.com", PhoneNumber: "9930627430"},
	}

	// Insert the users into the database
	result := db.Create(&users)
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	fmt.Println("Users inserted successfully")

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

		isBooked, err := strconv.ParseBool(record[3])
		if err != nil {
			log.Printf("Invalid IsBooked value for record: %v", record)
			continue
		}

		seat := models.Seat{
			ID:             int(id),
			SeatIdentifier: record[1],
			SeatClass:      record[2],
			IsBooked:       isBooked,
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

	file2, err := os.Open("seats_price.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file2.Close()

	// Create a new CSV reader for the second file
	reader2 := csv.NewReader(bufio.NewReader(file2))

	// Read the CSV records from the second file
	records2, err := reader2.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Iterate over the records from the second file and insert them into the database
	for _, record := range records2 {
		// Parse the record values
		id, err := strconv.Atoi(record[0])
		if err != nil {
			log.Printf("Invalid ID for record: %v", record)
			continue
		}

		minPrice, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			log.Printf("Invalid min price for record: %v", record)
			continue
		}

		normalPrice, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Printf("Invalid normal price for record: %v", record)
			continue
		}

		maxPrice, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			log.Printf("Invalid max price for record: %v", record)
			continue
		}

		seatPricing := models.SeatPricing{
			SeatPricingId: (id),
			SeatClass:     record[1],
			MinPrice:      minPrice,
			NormalPrice:   normalPrice,
			MaxPrice:      maxPrice,
		}

		// Insert the seat pricing into the database
		result := db.Create(&seatPricing)
		if result.Error != nil {
			log.Printf("Error inserting seat pricing: %v", result.Error)
			continue
		}

		log.Printf("Seat pricing inserted successfully: ID=%d", seatPricing.SeatPricingId)

	}
	DB = db
	return nil
}
func dropTables(db *gorm.DB) {
	db.Migrator().DropTable(&models.Seat{})
	db.Migrator().DropTable(&models.SeatPricing{})
	db.Migrator().DropTable(&models.Users{})
	db.Migrator().DropTable(&models.Bookings{})
}

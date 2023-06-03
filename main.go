package main

import (
	"encoding/csv"
	"flurn_assignment/controllers"
	"flurn_assignment/database"
	"flurn_assignment/models"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	r.GET("/seats", controllers.GetAllSeats)
	r.GET("/seats/:id", controllers.GetSeatPricing)
	r.POST("/booking", controllers.CreateBooking)
	r.GET("/bookings", controllers.GetBookings)

	r.GET("/export-seats-csv", func(c *gin.Context) {
		// Query the database to retrieve the data
		var seats []models.Seat
		if err := database.DB.Find(&seats).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to retrieve data from the database"})
			return
		}

		// Create a CSV file
		file, err := os.Create("data.seats.csv")
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create CSV file"})
			return
		}
		defer file.Close()

		// Create a CSV writer
		writer := csv.NewWriter(file)

		// Write the CSV header
		writer.Write([]string{"ID", "Seat Identifier", "Seat Class"})

		// Write each row of data to the CSV file
		for _, seat := range seats {
			writer.Write([]string{
				fmt.Sprint(seat.ID),
				seat.SeatIdentifier,
				seat.SeatClass,
			})
		}

		// Flush the CSV writer
		writer.Flush()

		// Check for any errors during CSV writing
		if err := writer.Error(); err != nil {
			c.JSON(500, gin.H{"error": "Failed to write data to CSV file"})
			return
		}

		c.JSON(200, gin.H{"message": "CSV file exported successfully"})
	})

	r.GET("/export-seats-pricing-csv", func(c *gin.Context) {
		// Query the database to retrieve the data
		var seats []models.Seat
		if err := database.DB.Find(&seats).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to retrieve data from the database"})
			return
		}

		// Create a CSV file
		file, err := os.Create("data.seats_price.csv")
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create CSV file"})
			return
		}
		defer file.Close()

		// Create a CSV writer
		writer := csv.NewWriter(file)

		// Write the CSV header
		writer.Write([]string{"ID", "Seat Identifier", "Seat Class"})

		// Write each row of data to the CSV file
		for _, seat := range seats {
			writer.Write([]string{
				fmt.Sprint(seat.ID),
				seat.SeatIdentifier,
				seat.SeatClass,
			})
		}

		// Flush the CSV writer
		writer.Flush()

		// Check for any errors during CSV writing
		if err := writer.Error(); err != nil {
			c.JSON(500, gin.H{"error": "Failed to write data to CSV file"})
			return
		}

		c.JSON(200, gin.H{"message": "CSV file exported successfully"})
	})

	err = r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

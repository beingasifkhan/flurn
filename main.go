package main

import (
	"flurn_assignment/controllers"
	"flurn_assignment/database"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	r := gin.Default()
	fmt.Println("Welcome to Flurn API's")

	r.GET("/seats", controllers.GetAllSeats)
	r.GET("/seats/:seat_identifier", controllers.GetSeatPricing)
	r.POST("/booking", controllers.CreateBooking)
	r.GET("/bookings", controllers.GetBookings)

	err = r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

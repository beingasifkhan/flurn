package main

import (
	"T_Booking_System/controllers"
	"T_Booking_System/database"
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
	fmt.Println("Welcome to Ticket bookings")

	r.GET("/seats", controllers.GetAllSeats)
	r.GET("/seats/:seat_identifier", controllers.GetSeatPricing)
	r.POST("/booking", controllers.CreateBooking)
	r.GET("/bookings", controllers.GetBookings)

	err = r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

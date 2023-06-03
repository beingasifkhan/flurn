package main

import (
	"flurn_assignment/controllers"
	"flurn_assignment/database"
	"log"

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

	err = r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

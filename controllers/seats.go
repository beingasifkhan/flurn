package controllers

import (
	"T_Booking_System/database"
	"T_Booking_System/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetAllSeats(c *gin.Context) {
	var seats []models.Seat
	err := database.DB.Find(&seats).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, seats)
}

func GetSeatPricing(c *gin.Context) {
	SeatIdentifier := c.Param("seat_identifier")

	var seat models.Seat
	err := database.DB.Where("seat_identifier = ?", SeatIdentifier).First(&seat).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var pricing models.SeatPricing
	err = database.DB.Where("seat_class = ?", seat.SeatClass).First(&pricing).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	percentage, err := CalculatePercentageBooked(seat.SeatClass)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var price float64
	if percentage < 40 {
		if pricing.MinPrice == 0 {
			price = pricing.NormalPrice
		} else {
			price = pricing.MinPrice
		}
	} else if percentage >= 40 && percentage <= 60 {
		if pricing.NormalPrice == 0 {
			price = pricing.MaxPrice
		} else {
			price = pricing.NormalPrice
		}
	} else {
		if pricing.MaxPrice == 0 {
			price = pricing.NormalPrice
		} else {
			price = pricing.MaxPrice
		}
	}

	// Return seat details along with pricing
	seatDetails := models.SeatDetails{
		SeatID:         seat.ID,
		SeatIdentifier: SeatIdentifier,
		SeatClass:      seat.SeatClass,
		Price:          price,
		Percentage:     percentage,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	c.JSON(http.StatusOK, seatDetails)
}

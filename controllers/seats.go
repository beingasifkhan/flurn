package controllers

import (
	"flurn_assignment/database"
	"flurn_assignment/models"
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

	for i := range seats {
		seats[i].IsBooked = seats[i].IsBooked
	}

	c.JSON(http.StatusOK, seats)
}

func GetSeatPricing(c *gin.Context) {
	Id := c.Param("id")

	var seat models.Seat
	err := database.DB.Where("id = ?", Id).First(&seat).Error
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

	percentage, err := calculatePercentageBooked(seat.SeatClass)
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
	seatDetails := gin.H{
		"seat_id":         Id,
		"seat_identifier": seat.SeatIdentifier,
		"seat_class":      seat.SeatClass,
		"price":           price,
		"percentage":      percentage,
		"created_at":      time.Now(),
		"update_at":       time.Now(),
	}

	c.JSON(http.StatusOK, seatDetails)
}

package controllers

import (
	"flurn_assignment/database"
	"flurn_assignment/models"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateBooking(c *gin.Context) {
	var requestBody models.BookingReq

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Check if selected seats are available
	var seats []models.Seat
	err := database.DB.Find(&seats, requestBody.SeatIDs).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, seat := range seats {
		if seat.IsBooked {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Seat %d is already booked", seat.ID)})
			return
		}
	}

	// Calculate total amount
	var totalAmount float64
	for _, seat := range seats {
		var pricing models.SeatPricing
		err := database.DB.Where("seat_class = ?", seat.SeatClass).First(&pricing).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Calculate percentage of seats booked
		percentage, err := CalculatePercentageBooked(seat.SeatClass)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Apply pricing rules based on the percentage
		var price float64
		if percentage < 40 {
			if pricing.MinPrice == 0 {
				price = pricing.NormalPrice
			} else {
				price = pricing.MinPrice
			}
		} else if percentage >= 40 && percentage <= 60 {
			price = pricing.NormalPrice
		} else {
			if pricing.MaxPrice == 0 {
				price = pricing.NormalPrice
			} else {
				price = pricing.MaxPrice
			}
		}

		totalAmount += price
	}
	totalAmount = math.Round(totalAmount*100) / 100
	formattedAmount := strconv.FormatFloat(totalAmount, 'f', 2, 64)

	var user models.Users
	err = database.DB.Where(" phone_number = ?", requestBody.PhoneNumber).First(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = database.DB.Where(" user_name = ?", requestBody.UserName).First(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate booking PID (unique identifier)
	bookingPID := GeneratePID()

	// Create new booking record
	booking := models.Bookings{
		BookingsPid:   bookingPID,
		SeatIDs:       requestBody.SeatIDs,
		UserName:      requestBody.UserName,
		PhoneNumber:   requestBody.PhoneNumber,
		UserID:        user.UserID,
		BookingAmount: totalAmount,
	}

	err = database.DB.Create(&booking).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = database.DB.Model(&models.Seat{}).Where("id IN (?)", requestBody.SeatIDs).Update("is_booked", true).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Booking created, but seat status update failed"})
	}

	// Return booking details
	bookingDetails := models.BookingDetails{
		UserName:    requestBody.UserName,
		PhoneNumber: requestBody.PhoneNumber,
		BookingID:   booking.BookingsId,
		BookingPID:  booking.BookingsPid,
		SeatIDs:     booking.SeatIDs,
		Amount:      formattedAmount,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	c.JSON(http.StatusOK, bookingDetails)
}

func GetBookings(c *gin.Context) {
	userIdentifier := c.Query("userIdentifier")

	if userIdentifier == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No user identifier provided"})
		return
	}

	var user models.Users
	if err := database.DB.Where("email = ? OR phone_number = ?", userIdentifier, userIdentifier).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	//Retrieve the bookings for the user
	var bookings []models.Bookings
	if err := database.DB.Where("user_id = ?", user.UserID).Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve bookings"})
		fmt.Println("Error retrieving bookings:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"bookings": bookings})
}

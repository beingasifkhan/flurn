package controllers

import (
	"flurn_assignment/database"
	"flurn_assignment/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllSeats(c *gin.Context) {
	var seats []models.Seat
	result := database.DB.Find(&seats)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"seats": seats})
}

func GetSeatPricing(c *gin.Context) {
	seatID := c.Param("id")

	var pricing models.SeatPricing
	result := database.DB.First(&pricing, seatID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Seat pricing not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}

	bookingsPercentage := float64(pricing.Bookings) / float64(maxBookingsForSeatClass(pricing.SeatClass))
	var price float64
	switch {
	case bookingsPercentage < 0.4:
		if pricing.MinPrice != 0 {
			price = pricing.MinPrice
		} else {
			price = pricing.NormalPrice
		}
	case bookingsPercentage >= 0.4 && bookingsPercentage <= 0.6:
		if pricing.NormalPrice != 0 {
			price = pricing.NormalPrice
		} else {
			price = pricing.MaxPrice
		}
	default:
		if pricing.MaxPrice != 0 {
			price = pricing.MaxPrice
		} else {
			price = pricing.NormalPrice
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"seat":  seatID,
		"price": price,
	})
}
func maxBookingsForSeatClass(seatClass string) int {

	switch seatClass {
	case "FirstClass":
		return 35
	case "BusinessClass":
		return 25
	case "EconomyClass":
		return 40
	default:
		return 0
	}
}
func CreateBooking(c *gin.Context) {
	var bookingReq models.BookingRequest

	err := c.ShouldBindJSON(&bookingReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if any of the requested seats are already booked
	var bookedSeats []models.Seat
	result := database.DB.
		Where("id IN ?", bookingReq.SeatIDs).
		Where("is_booked = ?", true).
		Find(&bookedSeats)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if len(bookedSeats) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Some seats are already booked"})
		return
	}

	// Update the booked status of the requested seats
	result = database.DB.Model(&models.Seat{}).
		Where("id IN ?", bookingReq.SeatIDs).
		Update("is_booked", true)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Create a new booking record
	booking := models.Booking{
		Name:        bookingReq.Name,
		PhoneNumber: bookingReq.PhoneNumber,
		SeatIDs:     bookingReq.SeatIDs,
	}
	result = database.DB.Create(&booking)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	totalAmount, err := calculateTotalAmount(booking.SeatIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"booking_id":   booking.ID,
		"total_amount": totalAmount,
	})

}

func GetBookings(c *gin.Context) {
	userIdentifier := c.Query("userIdentifier")
	if userIdentifier == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User identifier is required"})
		return
	}

	var bookings []models.Booking
	result := database.DB.
		Preload("Seats").
		Where("name = ? OR phone_number = ?", userIdentifier, userIdentifier).
		Find(&bookings)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"bookings": bookings})
}

func calculateTotalAmount(seatIDs []uint) (float64, error) {
	totalAmount := 0.0

	for _, seatID := range seatIDs {
		seatPricing := models.SeatPricing{}
		result := database.DB.Where("seat_id = ?", seatID).First(&seatPricing)
		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			return 0, result.Error
		}

		bookingsCount := seatPricing.BookingsCount
		bookingsPercentage := float64(bookingsCount) / float64(maxBookingsForSeatClass(seatPricing.SeatClass)) * 100

		var price float64

		if bookingsPercentage >= 60 {
			price = seatPricing.MaxPrice
		} else if bookingsPercentage >= 40 {
			price = seatPricing.NormalPrice
		} else {
			price = seatPricing.MinPrice
		}

		totalAmount += price
	}

	return totalAmount, nil
}

package controllers

import (
	"T_Booking_System/database"
	"T_Booking_System/models"
	"crypto/rand"
	"encoding/hex"
	"strconv"
)

func CalculatePercentageBooked(seatClass string) (float64, error) {
	var totalSeats int64
	var bookedSeats int64
	err := database.DB.Model(&models.Seat{}).Where("seat_class = ?", seatClass).Count(&totalSeats).Error
	if err != nil {
		return 0, err
	}

	err = database.DB.Model(&models.Seat{}).Where("seat_class = ? AND is_booked = ?", seatClass, true).Count(&bookedSeats).Error
	if err != nil {
		return 0, err
	}

	percentage := (float64(bookedSeats) / float64(totalSeats)) * 100
	formattedPercentage, err := strconv.ParseFloat(strconv.FormatFloat(percentage, 'f', 2, 64), 64)
	if err != nil {
		return 0, err
	}

	return formattedPercentage, nil
}

func GeneratePID() string {
	prefix := "flurn_"
	randomBytes := make([]byte, 4)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	pid := prefix + hex.EncodeToString(randomBytes)[:8]
	return pid
}

package database

import (
	"flurn_assignment/models"

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
	db.AutoMigrate(&models.Seat{}, &models.SeatPricing{}, &models.Booking{})

	DB = db
	return nil
}

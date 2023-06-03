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

	err = db.AutoMigrate(&models.Seat{})
	if err != nil {
		panic("Failed to migrate the Seat struct: " + err.Error())
	}
	err = db.AutoMigrate(&models.Booking{})
	if err != nil {
		panic("failed to create table")
	}
	err = db.AutoMigrate(&models.SeatPricing{})
	if err != nil {
		panic("failed to create table")
	}

	DB = db
	return nil
}

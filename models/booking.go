package models

import "time"

type Users struct {
	UserID      int    `json:"user_id" gorm:"primaryKey"`
	UserName    string `json:"user_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

type Bookings struct {
	UserName      string  `json:"user_name"`
	PhoneNumber   string  `json:"phone_number"`
	BookingsId    int     `json:"bookings_id" gorm:"primaryKey"`
	BookingsPid   string  `json:"bookings_pid"`
	SeatIDs       []uint  `json:"seat_ids" gorm:"-"`
	UserID        int     `json:"user_id" gorm:"foreignKey:UserID"`
	BookingAmount float64 `json:"booking_amount"`
}

type BookingsData struct {
	BookingsId    int
	BookingsPid   string
	SeatsID       []int
	UsersID       int
	BookingAmount float64
}

type BookingReq struct {
	SeatIDs     []uint `json:"seat_ids"`
	UserName    string `json:"user_name"`
	PhoneNumber string `json:"phone_number"`
}
type BookingDetails struct {
	UserName    string    `json:"user_name"`
	PhoneNumber string    `json:"phone_number"`
	BookingID   int       `json:"booking_id"`
	BookingPID  string    `json:"booking_pid"`
	SeatIDs     []uint    `json:"seat_ids"`
	Amount      string    `json:"amount"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

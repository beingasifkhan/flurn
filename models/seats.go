package models

type Seat struct {
	ID             uint   `json:"id" gorm:"primaryKey"`
	SeatIdentifier string `json:"seat_identifier"`
	SeatClass      string `json:"seat_class"`
}

type SeatPricing struct {
	ID            uint    `json:"id" gorm:"primaryKey"`
	SeatClass     string  `json:"seat_class"`
	MinPrice      float64 `json:"min_price"`
	NormalPrice   float64 `json:"normal_price"`
	MaxPrice      float64 `json:"max_price"`
	Bookings      int     `json:"bookings"`
	BookingsCount uint    `json:"bookings_count" gorm:"-"`
}
type BookingRequest struct {
	SeatIDs     []uint `json:"seat_ids"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}
type Booking struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	SeatIDs     []uint `json:"seat_ids" gorm:"-"`
	Seats       []Seat `json:"seats" gorm:"many2many:booking_seats"`
}

type BookingSeat struct {
	BookingID uint `gorm:"primaryKey"`
	SeatID    uint `gorm:"primaryKey"`
}

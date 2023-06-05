package models

type Seat struct {
	ID             int    `json:"id" gorm:"primaryKey"`
	SeatIdentifier string `json:"seat_identifier"`
	SeatClass      string `json:"seat_class"`
	IsBooked       bool   `json:"is_booked"`
}

type SeatPricing struct {
	SeatPricingId int     `json:"seat_pricing_id" gorm:"primaryKey"`
	SeatClass     string  `json:"seat_class"`
	MinPrice      float64 `json:"min_price"`
	NormalPrice   float64 `json:"normal_price"`
	MaxPrice      float64 `json:"max_price"`
}

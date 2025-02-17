package entity

type Room struct {
	ID           string  `json:"id"`
	Type         string  `json:"type"`         // room_type (Enum: e.g., "single", "double", etc.)
	Category     string  `json:"category"`     // room_category (Enum: e.g., "standard", "deluxe", etc.)
	Status       string  `json:"status"`       // room_status (Enum: e.g., "available", "occupied", etc.)
	Price        float64 `json:"price"`        // Decimal value
	Availability bool    `json:"availability"` // True (available) or False (unavailable)
	Rating       float64 `json:"rating"`       // Average rating
	CreatedAt    string  `json:"created_at"`   // Timestamp
	UpdatedAt    string  `json:"updated_at"`   // Timestamp
}

type RoomList struct {
	Items []Room `json:"rooms"`
	Count int    `json:"count"`
}

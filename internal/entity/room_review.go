package entity

type RoomReview struct {
	ID        string  `json:"id"`         // UUID
	UserID    string  `json:"user_id"`    // Foydalanuvchi identifikatori (UUID)
	RoomID    string  `json:"room_id"`    // Xona identifikatori (UUID)
	Rating    float64 `json:"rating"`     // Reyting qiymati
	Comment   string  `json:"comment"`    // Fikr-mulohaza
	CreatedAt string  `json:"created_at"` // Yaratilgan vaqt (timestamp)
	UpdatedAt string  `json:"updated_at"` // So'nggi yangilanish vaqti (timestamp)
}

type RoomReviewList struct {
	Items []RoomReview `json:"room_reviews"`
	Count int          `json:"count"`
}

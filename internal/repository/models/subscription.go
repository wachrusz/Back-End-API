package models

// Subscription представляет собой информацию о подписке пользователя.
type Subscription struct {
	ID        int64  `json:"id"`
	UserID    string `json:"user_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	IsActive  bool   `json:"is_active"`
}

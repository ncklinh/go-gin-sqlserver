package models

type CartSummary struct {
	UserId     int        `json:"user_id"`
	TotalPrice float64    `json:"total_price"`
	TotalItems int        `json:"total_items"`
	Items      []CartItem `json:"items"`
}

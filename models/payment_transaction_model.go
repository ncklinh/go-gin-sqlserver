package models

import "time"

type PaymentTransaction struct {
	Id         int       `json:"id" gorm:"primaryKey"`
	MethodCode string    `json:"method_code" gorm:"not null"`
	OrderId    string    `json:"order_id"`
	Amount     float64   `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
}

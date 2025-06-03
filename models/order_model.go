package models

type Order struct {
	Id            int           `json:"id" gorm:"primaryKey"`
	Items         []OrderItem   `json:"items" gorm:"foreignKey:OrderId"`
	TotalAmount   float64       `json:"total_amount"`
	PaymentStatus PaymentStatus `json:"payment_status"`
}

type PaymentStatus string

const (
	Paid           PaymentStatus = "Paid"
	PaymentWaiting PaymentStatus = "PaymentWaiting"
)

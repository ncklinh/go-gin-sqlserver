package models

type OrderItem struct {
	Id        int     `json:"id" gorm:"primaryKey"`
	OrderId   int     `json:"order_id"`

	ProductId int     `json:"product_id"`
	Product   Product `json:"product" gorm:"foreignKey:ProductId"`

	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

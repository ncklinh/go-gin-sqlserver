package models

type CartItem struct {
	Id int `json:"id" gorm:"primaryKey;autoIncrement:true;unique"`
	UserId int `json:"user_id"`
	ProductId int `json:"product_id"`
	Quantity int `json:"quantity"`
	Product Product `json:"product" gorm:"foreignKey:ProductId"`
}
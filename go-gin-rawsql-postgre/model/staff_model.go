package model

import "time"

type Staff struct {
	StaffId    int       `json:"staff_id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	AddressId  int       `json:"address_id"`
	Email      string    `json:"email"`
	StoreId    string    `json:"store_id"`
	Active     bool      `json:"active"`
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	LastUpdate time.Time `json:"last_update"`
	Picture    []byte    `json:"picture"`
}

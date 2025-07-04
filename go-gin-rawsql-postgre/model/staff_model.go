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
	Password   string    `json:"-"`
	Role       string    `json:"role"`
	LastUpdate time.Time `json:"last_update"`
	Picture    []byte    `json:"picture"`
}

// CreateStaffRequest is used for creating new staff members
type CreateStaffRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	AddressId int    `json:"address_id"`
	Email     string `json:"email" binding:"required"`
	StoreId   string `json:"store_id"`
	Active    bool   `json:"active"`
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Role      string `json:"role"`
	Picture   []byte `json:"picture"`
}

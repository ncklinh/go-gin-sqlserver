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
	Picture   []byte `json:"picture"`
}

// LoginRequest is used for staff login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Package model represent objects structure in application
package model

import "encoding/json"

// Order type represent order structure in database
type Order struct {
	OrderID     string `json:"orderID"`
	OrderName   string `json:"orderName"`
	OrderCost   int    `json:"orderCost"`
	IsDelivered bool   `json:"isDelivered"`
}

// AuthUser struct represents user information
type AuthUser struct {
	UserUUID     string `json:"userID"`
	UserName     string `json:"userName"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
}

func (order Order) MarshalBinary() ([]byte, error) {
	return json.Marshal(order)
}

func (order Order) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &order)
}

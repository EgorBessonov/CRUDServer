// Package model represent objects structure in application
package model

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

// Cat struct represents cat information
type Cat struct{
	CatName string
	CatAge int
	IsHungry bool
}
package repository

import(
	"github.com/google/uuid"
)

type User struct{
	UserId uuid.UUID
	UserName string	`json:"userName" bson:"userName"`
	UserAge int `json:"userAge" bson:"userAge"`
	IsAdult bool `json:"isAdult" bson:"isAdult"`
}


type IRepository interface{
	CreateUser(u User)error
	ReadUser()
	UpdateUser(u User)error
	DeleteUser(userID uuid.UUID)error
	AddImage()
	GetImage()
}

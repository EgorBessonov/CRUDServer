package repository


type User struct{
	UserId string `json: "_id" bson: "_id"`
	UserName string	`json:"userName" bson:"userName"`
	UserAge int `json:"userAge" bson:"userAge"`
	IsAdult bool `json:"isAdult" bson:"isAdult"`
}


type IRepository interface{
	CreateUser(u User)error
	ReadUser(u string) (*User, error)
	UpdateUser(u User)error
	DeleteUser(userID string)error
	AddImage()
	GetImage()
}

// Package repository replies for database access
package repository

// Config type replies for connection to current database
type Config struct {
	CurrentDB     string `env:"CURRENTDB" envDefault:"postgres"`
	PostgresdbUrl string `env:"POSTGRESDB_URL"`
	MongodbUrl    string `env:"MONGODB_URL"`
}

// User type represent user structure in database
type User struct {
	UserID   string
	UserName string `json:"userName" bson:"userName"`
	UserAge  int    `json:"userAge" bson:"userAge"`
	IsAdult  bool   `json:"isAdult" bson:"isAdult"`
}

// IRepository interface represent repository behavior
type IRepository interface {
	CreateUser(u User) error
	ReadUser(u string) (User, error)
	UpdateUser(u User) error
	DeleteUser(userID string) error
	AddImage()
	GetImage()
}

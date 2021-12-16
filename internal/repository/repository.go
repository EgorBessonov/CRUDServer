// Package repository replies for database access
package repository

import(
	"time"
)

// Config type replies for connection to current database
type Config struct {
	CurrentDB     string `env:"CURRENTDB" envDefault:"postgres"`
	PostgresdbURL string `env:"POSTGRESDB_URL"`
	MongodbURL    string `env:"MONGODB_URL"`
}

// RegistrationForm struct represents user information
type RegistrationForm struct {
	UserID       int    `json:"userID" bson:"userID"`
	UserName     string `json:"userName" bson:"userName"`
	Email        string `json:"email" bson:"email"`
	Password     string `json:"password" bson:"password"`
	RefreshToken string `json:"refreshToken" bson:"refreshToken"`
	ExpiresIn string `json:"expiresIn" bson:"expiresIn"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
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
	CreateUser(User) error
	ReadUser(string) (User, error)
	UpdateUser(User) error
	DeleteUser(string) error
	AddImage()
	GetImage()
	CreateAuthUser(RegistrationForm) error
	GetAuthUser(string) (RegistrationForm, error)
	CloseDBConnection() error
}

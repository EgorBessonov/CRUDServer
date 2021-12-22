// Package repository replies for database access
package repository

import "context"

// Config type replies for connection to current database
type Config struct {
	SecretKey     string `env:"SECRETKEY"`
	CurrentDB     string `env:"CURRENTDB" envDefault:"postgres"`
	PostgresdbURL string `env:"POSTGRESDB_URL"`
	MongodbURL    string `env:"MONGODB_URL"`
}

// RegistrationForm struct represents user information
type RegistrationForm struct {
	UserUUID     string `json:"userID"`
	UserName     string `json:"userName"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
}

// User type represent user structure in database
type User struct {
	UserID   string
	UserName string `json:"userName" bson:"userName"`
	UserAge  int    `json:"userAge" bson:"userAge"`
	IsAdult  bool   `json:"isAdult" bson:"isAdult"`
}

// Repository interface represent repository behavior
type Repository interface {
	CreateUser(User, context.Context) error
	ReadUser(string, context.Context) (User, error)
	UpdateUser(User, context.Context) error
	DeleteUser(string, context.Context) error
	CreateAuthUser(RegistrationForm, context.Context) error
	GetAuthUser(string, context.Context) (RegistrationForm, error)
	GetAuthUserByID(string, context.Context) (RegistrationForm, error)
	UpdateAuthUser(email, refreshToken string, ctx context.Context) error
	CloseDBConnection() error
}

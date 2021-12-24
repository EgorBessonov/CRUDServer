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
	CreateUser(context.Context, User) error
	ReadUser(context.Context, string) (User, error)
	UpdateUser(context.Context, User) error
	DeleteUser(context.Context, string) error
	CreateAuthUser(context.Context, *RegistrationForm) error
	GetAuthUser(context.Context, string) (RegistrationForm, error)
	GetAuthUserByID(context.Context, string) (RegistrationForm, error)
	UpdateAuthUser(ctx context.Context, email, refreshToken string) error
	CloseDBConnection() error
}

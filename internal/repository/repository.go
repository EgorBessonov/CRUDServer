// Package repository replies for database access
package repository

// Config type replies for connection to current database
type Config struct {
	CurrentDB     string `env:"CURRENTDB" envDefault:"postgres"`
	PostgresdbURL string `env:"POSTGRESDB_URL"`
	MongodbURL    string `env:"MONGODB_URL"`
}

// AuthForm struct represents user login information
type AuthForm struct {
	UserID       int    `json:"userID" bson:"userID"`
	UserName     string `json:"userName" bson:"userName"`
	Email        string `json:"email" bson:"email"`
	Password     string `json:"password" bson:"password"`
	RefreshToken string `json:"refreshToken" bson:"refreshToken"`
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
	CreateAuthUser(AuthForm) error
	GetAuthUser(string) (AuthForm, error)
	CloseDBConnection() error
}

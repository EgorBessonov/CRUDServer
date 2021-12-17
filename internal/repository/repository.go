// Package repository replies for database access
package repository

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

// IRepository interface represent repository behavior
type IRepository interface {
	CreateUser(User) error
	ReadUser(string) (User, error)
	UpdateUser(User) error
	DeleteUser(string) error
	CreateAuthUser(RegistrationForm) error
	GetAuthUser(string) (RegistrationForm, error)
	UpdateAuthUser(email string, refreshtoken string) error
	CloseDBConnection() error
}

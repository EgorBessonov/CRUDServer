// Package configs represents config structure
package configs

// Config type store all env info
type Config struct {
	SecretKey     string `env:"SECRETKEY"`
	CurrentDB     string `env:"CURRENTDB" envDefault:"postgres"`
	PostgresdbURL string `env:"POSTGRESDB_URL"`
	MongodbURL    string `env:"MONGODB_URL"`
	RedisURL      string `env:"REDISDB_URL"`
	StreamName    string `env:"STREAMNAME"`
	ServiceUUID   string `env:"SERVICEUUID"`
}

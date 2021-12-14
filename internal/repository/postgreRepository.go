package repository

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

// PostgresRepository type replies for accessing to postgres database
type PostgresRepository struct {
	DBconn *pgxpool.Pool
}

// CreateUser save User object into postgresql database
func (rps PostgresRepository) CreateUser(u User) error {
	result, err := rps.DBconn.Exec(context.Background(), "insert into users (userName, userAge, isAdult) "+
		"values ($1, $2, $3)", u.UserName, u.UserAge, u.IsAdult)
	if err != nil {
		postgresOperationError(err, "CreateUser()")
		return err
	}
	postgresOperationSuccess(result, "CreateUser()")
	return nil
}

// ReadUser returns User object from postgresql database
// with selection by UserId
func (rps PostgresRepository) ReadUser(u string) (User, error) {
	var user User
	err := rps.DBconn.QueryRow(context.Background(), "select userName, userAge, isAdult from users "+
		"where userId=$1", u).Scan(&user.UserName, &user.UserAge, &user.IsAdult)
	if err != nil {
		postgresOperationError(err, "ReadUser()")
		return User{}, err
	}
	postgresOperationSuccess(nil, "ReadUser()")
	return user, nil
}

// UpdateUser update User object from postgresql database
// with selection by UserId
func (rps PostgresRepository) UpdateUser(u User) error {
	result, err := rps.DBconn.Exec(context.Background(), "update users "+
		"set userName=$2, userAge=$3, isAdult=$4"+
		"where userid=$1", u.UserID, u.UserName, u.UserAge, u.IsAdult)

	if err != nil {
		postgresOperationError(err, "UpdateUser()")
		return err
	}
	postgresOperationSuccess(result, "UpdateUser()")
	return nil
}

// DeleteUser delete User object from postgresql database
// with selection by UserId
func (rps PostgresRepository) DeleteUser(userID string) error {
	result, err := rps.DBconn.Exec(context.Background(), "delete from users where userId=$1", userID)
	if err != nil {
		postgresOperationError(err, "DeleteUser()")
		return err
	}
	postgresOperationSuccess(result, "DeleteUser()")
	return nil
}

// AddImage function
func (rps PostgresRepository) AddImage() {

}

// GetImage function
func (rps PostgresRepository) GetImage() {

}

// CreateAuthUser save authentication info about user into
// postgres database
func (rps PostgresRepository) CreateAuthUser(lf AuthForm) error {
	result, err := rps.DBconn.Exec(context.Background(), "insert into authusers (username, email, password)"+
		"values($1, $2, $3)", lf.UserName, lf.Email, lf.Password)
	if err != nil {
		postgresOperationError(err, "CreateAuthUser()")
		return err
	}
	postgresOperationSuccess(result, "CreateAuthUser()")
	return nil
}

// GetAuthUser return authentication info about user into
// postgres database
func (rps PostgresRepository) GetAuthUser(string) (AuthForm, error) {
	return AuthForm{}, nil
}

// CloseDBConnection is using to close current postgres database connection
func (rps PostgresRepository) CloseDBConnection() error {
	rps.DBconn.Close()
	postgresOperationSuccess(nil, "CloseDBConnection")
	return nil
}

func postgresOperationError(err error, method string) {
	log.WithFields(log.Fields{
		"method": method,
		"status": "Operation failed.",
		"error":  err,
	}).Info("Postgres repository info.")
}

func postgresOperationSuccess(result pgconn.CommandTag, method string) {
	log.WithFields(log.Fields{
		"method": method,
		"status": "Operation ended successfully",
		"result": result,
	}).Info("Postgres repository info.")
}

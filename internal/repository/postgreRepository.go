package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

// PostgresRepository type replies for accessing to postgres database
type PostgresRepository struct {
	DBconn *pgxpool.Pool
}

// CreateUser save User object into postgresql database
func (rps PostgresRepository) CreateUser(ctx context.Context, u User) error {
	result, err := rps.DBconn.Exec(ctx, "insert into users (userName, userAge, isAdult) "+
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
func (rps PostgresRepository) ReadUser(ctx context.Context, u string) (User, error) {
	var user User
	err := rps.DBconn.QueryRow(ctx, "select userName, userAge, isAdult from users "+
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
func (rps PostgresRepository) UpdateUser(ctx context.Context, u User) error {
	result, err := rps.DBconn.Exec(ctx, "update users "+
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
func (rps PostgresRepository) DeleteUser(ctx context.Context, userID string) error {
	result, err := rps.DBconn.Exec(ctx, "delete from users where userId=$1", userID)
	if err != nil {
		postgresOperationError(err, "DeleteUser()")
		return err
	}
	postgresOperationSuccess(result, "DeleteUser()")
	return nil
}

// CreateAuthUser method saves authentication info about user into
// postgres database
func (rps PostgresRepository) CreateAuthUser(ctx context.Context, lf RegistrationForm) error {
	fmt.Println(lf)
	result, err := rps.DBconn.Exec(ctx, "insert into authusers (username, email, password)"+
		"values($1, $2, $3)", lf.UserName, lf.Email, lf.Password)
	if err != nil {
		postgresOperationError(err, "CreateAuthUser()")
		return err
	}
	postgresOperationSuccess(result, "CreateAuthUser()")
	return nil
}

// GetAuthUser method returns authentication info about user from
// postgres database with selection by email
func (rps PostgresRepository) GetAuthUser(ctx context.Context, email string) (RegistrationForm, error) {
	var authUser RegistrationForm
	err := rps.DBconn.QueryRow(ctx, "select useruuid, username, email, password from authusers "+
		"where email=$1", email).Scan(&authUser.UserUUID, &authUser.UserName, &authUser.Email, &authUser.Password)
	if err != nil {
		postgresOperationError(err, "GetAuthUser()")
		return RegistrationForm{}, err
	}
	postgresOperationSuccess(nil, "GetAuthUser()")
	return authUser, nil
}

// GetAuthUserByID method returns authentication info about user from
// postgres database with selection by id
func (rps PostgresRepository) GetAuthUserByID(ctx context.Context, userUUID string) (RegistrationForm, error) {
	var authUser RegistrationForm
	err := rps.DBconn.QueryRow(ctx, "select useruuid, username, email, password, refreshtoken from authusers "+
		"where useruuid=$1", userUUID).Scan(&authUser.UserUUID, &authUser.UserName, &authUser.Email, &authUser.Password, &authUser.RefreshToken)
	if err != nil {
		postgresOperationError(err, "GetAuthUserByID()")
		return RegistrationForm{}, err
	}
	postgresOperationSuccess(nil, "GetAuthUserByID()")
	return authUser, nil
}

// UpdateAuthUser is method to set refresh token into authuser info
func (rps PostgresRepository) UpdateAuthUser(ctx context.Context, email, refreshToken string) error {
	result, err := rps.DBconn.Exec(ctx, "update authusers "+
		"set refreshtoken=$2"+
		"where email=$1", email, refreshToken)
	if err != nil {
		postgresOperationError(err, "UpdateAuthUser()")
		return err
	}
	postgresOperationSuccess(result, "UpdateAuthUser()")
	return nil
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

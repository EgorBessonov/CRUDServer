package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
	"time"
)

type PostgresRepository struct{}

//CreateUser save User object into postgresql database
// and returns creation status 
func (rps PostgresRepository) CreateUser(u User) error {
	conn, err := pgx.Connect(context.Background(), "postgresql://postgres:passwd@localhost:5614/crudserverdb")
	if err != nil {
		log.WithFields(log.Fields{
			"method": "CreateUser()",
			"status": "Failed connection to crudserverdb.",
			"error":  err,
		}).Info("Postgresql repository info.")
		return err
	} else {
		log.WithFields(log.Fields{
			"time":   time.Now(),
			"method": "CreateUser()",
			"status": "Successfully connected to crudserverdb.",
		}).Info("Postgresql repository info.")
	}
	defer conn.Close(context.Background())

	result, err := conn.Exec(context.Background(), "insert into users (userName, userAge, isAdult) "+
		"values ($1, $2, $3)", u.UserName, u.UserAge, u.IsAdult)
	if err != nil {
		log.WithFields(log.Fields{
			"method": "CreateUser()",
			"status": "Failed while inserting.",
			"error":  err,
		}).Info("Postgresql repository info.")
		return err
	} else {
		log.WithFields(log.Fields{
			"method": "CreateUser()",
			"status": "Successfully inserted.",
			"result": result,
		}).Info("Postgresql repository info.")
	}
	return nil
}

//ReadUser return User object from postgresql database
//with selection by UserId
func (rps PostgresRepository) ReadUser(u string) (User, error) {
	conn, err := pgx.Connect(context.Background(), "postgresql://postgres:passwd@localhost:5614/crudserverdb")
	if err != nil {
		log.WithFields(log.Fields{
			"method": "ReadUser()",
			"status": "Failed connection to crudserverdb.",
			"error":  err,
		}).Info("Postgresql repository info.")
		return User{}, err
	} else {
		log.WithFields(log.Fields{
			"time":   time.Now(),
			"method": "ReadUser()",
			"status": "Successfully connected to crudserverdb.",
		}).Info("Postgresql repository info.")
	}
	defer conn.Close(context.Background())
	var user User
	err = conn.QueryRow(context.Background(), "select userName, userAge, isAdult from users "+
		"where userId=$1", u).Scan(&user.UserName, &user.UserAge, &user.IsAdult)
	if err != nil {
		log.WithFields(log.Fields{
			"method": "ReadUser()",
			"status": "Failed while reading.",
			"error":  err,
		}).Info("Postgresql repository info.")
		return User{}, err
	} else {
		log.WithFields(log.Fields{
			"method": "ReadUser()",
			"status": "Successfully read.",
		}).Info("Postgresql repository info.")
	}
	return user, nil
}

//UpdateUser update User object from postgresql database
//with selection by UserId
//and returns updating status
func (rps PostgresRepository) UpdateUser(u User) error {
	conn, err := pgx.Connect(context.Background(), "postgresql://postgres:passwd@localhost:5614/crudserverdb")
	if err != nil {
		log.WithFields(log.Fields{
			"method": "UpdateUser()",
			"status": "Failed connection to crudserverdb.",
			"error":  err,
		}).Info("Postgresql repository info.")
		return err
	} else {
		log.WithFields(log.Fields{
			"time":   time.Now(),
			"method": "UpdateUser()",
			"status": "Successfully connected to crudserverdb.",
		}).Info("Postgresql repository info.")
	}
	defer conn.Close(context.Background())

	result, err := conn.Exec(context.Background(), "update users "+
		"set userName=$2, userAge=$3, isAdult=$4"+
		"where userid=$1", u.UserId, u.UserName, u.UserAge, u.IsAdult)
	if err != nil {
		log.WithFields(log.Fields{
			"method": "UpdateUser()",
			"status": "Failed while updating.",
			"error":  err,
		}).Info("Postgresql repository info.")
		return err
	} else {
		log.WithFields(log.Fields{
			"method": "UpdateUser()",
			"status": "Successfully updated.",
			"result": result,
		}).Info("Postgresql repository info.")
	}
	return nil
}

//DeleteUser delete User object from postgresql database
//with selection by UserId
//and returns deletion status
func (rps PostgresRepository) DeleteUser(userID string) error {
	conn, err := pgx.Connect(context.Background(), "postgresql://postgres:passwd@localhost:5614/crudserverdb")
	if err != nil {
		log.WithFields(log.Fields{
			"method": "DeleteUser()",
			"status": "Failed connection to crudserverdb.",
			"error":  err,
		}).Info("Postgresql repository info.")
		return err
	} else {
		log.WithFields(log.Fields{
			"time":   time.Now(),
			"method": "DeleteUser()",
			"status": "Successfully connected to crudserverdb.",
		}).Info("Postgresql repository info.")
	}
	defer conn.Close(context.Background())

	result, err := conn.Exec(context.Background(), "delete from users where userId=$1", userID)
	if err != nil {
		log.WithFields(log.Fields{
			"method": "DeleteUser()",
			"status": "Failed while deleting.",
			"error":  err,
		}).Info("Postgresql repository info.")
		return err
	} else {
		log.WithFields(log.Fields{
			"method": "DeleteUser()",
			"status": "Successfully deleted.",
			"result": result,
		}).Info("Postgresql repository info.")
	}
	return nil
}


//AddImage function
func (rps PostgresRepository) AddImage() {

}
// GetImage function
func (rps PostgresRepository) GetImage() {

}

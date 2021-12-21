package repository

import (
	"context"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// MongoRepository type replies for accessing to mongo database
type MongoRepository struct {
	DBconn *mongo.Client
}

// CreateUser method saves User object into mongo database
func (rps MongoRepository) CreateUser(u User) error {
	col := rps.DBconn.Database("crudserver").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := col.InsertOne(ctx, bson.D{
		{Key: "_id", Value: uuid.New().String()},
		{Key: "userName", Value: u.UserName},
		{Key: "userAge", Value: u.UserAge},
		{Key: "isAdult", Value: u.IsAdult},
	})
	if err != nil {
		mongoOperationError(err, "CreateUser()")
		return err
	}
	mongoOperationSuccess("CreateUser()")
	return nil
}

// ReadUser method returns User object from mongo database
// with selection by UserId
func (rps MongoRepository) ReadUser(userID string) (User, error) {
	col := rps.DBconn.Database("crudserver").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var rUser User
	err := col.FindOne(ctx, bson.D{{Key: "_id", Value: userID}}).Decode(&rUser)
	if err != nil {
		mongoOperationError(err, "ReadUser()")
		return User{}, err
	}
	mongoOperationSuccess("ReadUser()")
	return rUser, nil
}

// UpdateUser method updates User object from mongo database
// with selection by UserId
func (rps MongoRepository) UpdateUser(u User) error {
	col := rps.DBconn.Database("crudserver").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := col.UpdateOne(ctx, bson.D{{Key: "_id", Value: u.UserID}}, bson.D{
		{Key: "userName", Value: u.UserName},
		{Key: "userAge", Value: u.UserAge},
		{Key: "isAdult", Value: u.IsAdult},
	})
	if err != nil {
		mongoOperationError(err, "UpdateUser()")
		return err
	}
	mongoOperationSuccess("UpdateUser()")
	return nil
}

// DeleteUser method deletes User object from mongo database
// with selection by UserId
func (rps MongoRepository) DeleteUser(userID string) error {
	col := rps.DBconn.Database("crudserver").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := col.DeleteOne(ctx, bson.D{{Key: "_id", Value: userID}})
	if err != nil {
		mongoOperationError(err, "DeleteUser()")
		return err
	}
	mongoOperationSuccess("DeleteUser()")
	return nil
}

// GetAuthUser method returns authentication info about user from
// mongo database with selection by email
func (rps MongoRepository) GetAuthUser(email string) (RegistrationForm, error) {
	return RegistrationForm{}, nil
}

// GetAuthUserByID method returns authentication info about user from
// mongo database with selection by ID
func (rps MongoRepository) GetAuthUserByID(userUUID string) (RegistrationForm, error) {
	return RegistrationForm{}, nil
}

// CreateAuthUser method saves authentication info about user into
// postgres database
func (rps MongoRepository) CreateAuthUser(lf RegistrationForm) error {
	col := rps.DBconn.Database("crudserver").Collection("authusers")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := col.InsertOne(ctx, bson.D{
		{Key: "email", Value: lf.Email},
		{Key: "password", Value: lf.Password},
	})
	if err != nil {
		mongoOperationError(err, "CreateAuthUser()")
		return err
	}
	mongoOperationSuccess("CreateAuthUser()")
	return nil
}

// UpdateAuthUser method changes user refresh token
func (rps MongoRepository) UpdateAuthUser(email, refreshToken string) error {
	return nil
}

// CloseDBConnection is using for closing current mongo database connection
func (rps MongoRepository) CloseDBConnection() error {
	err := rps.DBconn.Disconnect(context.Background())
	if err != nil {
		mongoOperationError(err, "CloseDBConnection()")
	} else {
		mongoOperationSuccess("CloseDBConnection()")
	}
	return nil
}

func mongoOperationError(err error, method string) {
	log.WithFields(log.Fields{
		"method": method,
		"status": "Operation failed.",
		"error":  err,
	}).Info("Mongo repository info.")
}

func mongoOperationSuccess(method string) {
	log.WithFields(log.Fields{
		"method": method,
		"status": "Operation ended successfully",
	}).Info("Mongo repository info.")
}

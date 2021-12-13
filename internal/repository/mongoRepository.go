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

// CreateUser save User object into mongo database
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

// ReadUser returns User object from mongo database
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

// UpdateUser update User object from mongo database
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

// DeleteUser delete User object from mongo database
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

// AddImage function
func (rps MongoRepository) AddImage() {
	// TODO
}

// GetImage function
func (rps MongoRepository) GetImage() {
	// TODO
}

// GetAuthUser return authentication info about user into
// postgres database
func (rps MongoRepository) GetAuthUser(email string) (LoginForm, error) {
	return LoginForm{}, nil
}

// CreateAuthUser save authentication info about user into
// postgres database
func (rps MongoRepository) CreateAuthUser(LoginForm) error {

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

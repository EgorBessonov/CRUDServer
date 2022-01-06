package repository

import (
	"CRUDServer/internal/model"
	"context"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	timeout = 10
)

// MongoRepository type replies for accessing to mongo database
type MongoRepository struct {
	DBconn *mongo.Client
}

// Create method saves User object into mongo database
func (rps MongoRepository) Create(ctx context.Context, u model.Order) error {
	col := rps.DBconn.Database("crudserver").Collection("users")
	ctx, cancel := context.WithTimeout(ctx, timeout*time.Second)
	defer cancel()
	_, err := col.InsertOne(ctx, bson.D{
		{Key: "_id", Value: uuid.New().String()},
		{Key: "userName", Value: u.OrderName},
		{Key: "userAge", Value: u.OrderCost},
		{Key: "isAdult", Value: u.IsDelivered},
	})
	if err != nil {
		mongoOperationError(err, "Create()")
		return err
	}
	mongoOperationSuccess("Create()")
	return nil
}

// Read method returns User object from mongo database
// with selection by UserId
func (rps MongoRepository) Read(ctx context.Context, userID string) (model.Order, error) {
	col := rps.DBconn.Database("crudserver").Collection("users")
	ctx, cancel := context.WithTimeout(ctx, timeout*time.Second)
	defer cancel()

	var rUser model.Order
	err := col.FindOne(ctx, bson.D{{Key: "_id", Value: userID}}).Decode(&rUser)
	if err != nil {
		mongoOperationError(err, "Read()")
		return model.Order{}, err
	}
	mongoOperationSuccess("Read()")
	return rUser, nil
}

// Update method updates User object from mongo database
// with selection by UserId
func (rps MongoRepository) Update(ctx context.Context, u model.Order) error {
	col := rps.DBconn.Database("crudserver").Collection("users")
	ctx, cancel := context.WithTimeout(ctx, timeout*time.Second)
	defer cancel()

	_, err := col.UpdateOne(ctx, bson.D{{Key: "_id", Value: u.OrderID}}, bson.D{
		{Key: "userName", Value: u.OrderName},
		{Key: "userAge", Value: u.OrderCost},
		{Key: "isAdult", Value: u.IsDelivered},
	})
	if err != nil {
		mongoOperationError(err, "Update()")
		return err
	}
	mongoOperationSuccess("Update()")
	return nil
}

// Delete method deletes User object from mongo database
// with selection by UserId
func (rps MongoRepository) Delete(ctx context.Context, userID string) error {
	col := rps.DBconn.Database("crudserver").Collection("users")
	ctx, cancel := context.WithTimeout(ctx, timeout*time.Second)
	defer cancel()

	_, err := col.DeleteOne(ctx, bson.D{{Key: "_id", Value: userID}})
	if err != nil {
		mongoOperationError(err, "Delete()")
		return err
	}
	mongoOperationSuccess("Delete()")
	return nil
}

// GetAuthUser method returns authentication info about user from
// mongo database with selection by email
func (rps MongoRepository) GetAuthUser(ctx context.Context, email string) (model.AuthUser, error) {
	return model.AuthUser{}, nil
}

// GetAuthUserByID method returns authentication info about user from
// mongo database with selection by ID
func (rps MongoRepository) GetAuthUserByID(ctx context.Context, userUUID string) (model.AuthUser, error) {
	return model.AuthUser{}, nil
}

// CreateAuthUser method saves authentication info about user into
// postgres database
func (rps MongoRepository) CreateAuthUser(ctx context.Context, lf *model.AuthUser) error {
	col := rps.DBconn.Database("crudserver").Collection("authusers")
	ctx, cancel := context.WithTimeout(ctx, timeout*time.Second)
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
func (rps MongoRepository) UpdateAuthUser(ctx context.Context, email, refreshToken string) error {
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

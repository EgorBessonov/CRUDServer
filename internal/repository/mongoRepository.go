package repository

import (
	"context"
	"fmt"
	"github.com/EgorBessonov/CRUDServer/internal/model"
	"time"

	"github.com/google/uuid"
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
func (rps MongoRepository) Save(ctx context.Context, order *model.Order) error {
	col := rps.DBconn.Database("crudserver").Collection("users")
	ctx, cancel := context.WithTimeout(ctx, timeout*time.Second)
	defer cancel()
	_, err := col.InsertOne(ctx, bson.D{
		{Key: "_id", Value: uuid.New().String()},
		{Key: "userName", Value: order.OrderName},
		{Key: "userAge", Value: order.OrderCost},
		{Key: "isAdult", Value: order.IsDelivered},
	})
	if err != nil {
		return fmt.Errorf("mongo repository: can't save order - %w", err)
	}
	return nil
}

// Read method returns User object from mongo database
// with selection by UserId
func (rps MongoRepository) Get(ctx context.Context, orderID string) (*model.Order, error) {
	col := rps.DBconn.Database("crudserver").Collection("users")
	ctx, cancel := context.WithTimeout(ctx, timeout*time.Second)
	defer cancel()
	var order model.Order
	err := col.FindOne(ctx, bson.D{{Key: "_id", Value: orderID}}).Decode(&order)
	if err != nil {
		return nil, fmt.Errorf("mongo repository: can't get order - %w", err)
	}
	return &order, nil
}

// Update method updates User object from mongo database
// with selection by UserId
func (rps MongoRepository) Update(ctx context.Context, order *model.Order) error {
	col := rps.DBconn.Database("crudserver").Collection("users")
	ctx, cancel := context.WithTimeout(ctx, timeout*time.Second)
	defer cancel()
	_, err := col.UpdateOne(ctx, bson.D{{Key: "_id", Value: order.OrderID}}, bson.D{
		{Key: "userName", Value: order.OrderName},
		{Key: "userAge", Value: order.OrderCost},
		{Key: "isAdult", Value: order.IsDelivered},
	})
	if err != nil {
		return fmt.Errorf("mongo repository: can't update order - %w", err)
	}
	return nil
}

// Delete method deletes User object from mongo database
// with selection by UserId
func (rps MongoRepository) Delete(ctx context.Context, orderID string) error {
	col := rps.DBconn.Database("crudserver").Collection("users")
	ctx, cancel := context.WithTimeout(ctx, timeout*time.Second)
	defer cancel()
	_, err := col.DeleteOne(ctx, bson.D{{Key: "_id", Value: orderID}})
	if err != nil {
		return fmt.Errorf("mongo repository: can't delete order - %w", err)
	}
	return nil
}

// GetAuthUser method returns authentication info about user from
// mongo database with selection by email
func (rps MongoRepository) GetAuthUser(ctx context.Context, email string) (*model.AuthUser, error) {
	return nil, nil
}

// GetAuthUserByID method returns authentication info about user from
// mongo database with selection by ID
func (rps MongoRepository) GetAuthUserByID(ctx context.Context, userID string) (*model.AuthUser, error) {
	return nil, nil
}

// CreateAuthUser method saves authentication info about user into
// postgres database
func (rps MongoRepository) SaveAuthUser(ctx context.Context, authUser *model.AuthUser) error {
	col := rps.DBconn.Database("crudserver").Collection("authusers")
	ctx, cancel := context.WithTimeout(ctx, timeout*time.Second)
	defer cancel()

	_, err := col.InsertOne(ctx, bson.D{
		{Key: "email", Value: authUser.Email},
		{Key: "password", Value: authUser.Password},
	})
	if err != nil {
		return fmt.Errorf("mongo repository: can't save authUser - %w", err)
	}
	return nil
}

// UpdateAuthUser method changes user refresh token
func (rps MongoRepository) UpdateAuthUser(ctx context.Context, email, refreshToken string) error {
	return nil
}

// CloseDBConnection is using for closing current mongo database connection
func (rps MongoRepository) CloseDBConnection() error {
	if err := rps.DBconn.Disconnect(context.Background()); err != nil {
		return fmt.Errorf("mongo repository: database connection closed")
	}
	return nil
}

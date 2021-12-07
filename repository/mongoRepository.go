package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// MongoRepository type replies for accessing to mongo database
type MongoRepository struct{}

// CreateUser save User object into mongo database
func (rps MongoRepository) CreateUser(u User) error {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.WithFields(log.Fields{
			"method": "CreateUser()",
			"status": "Failed connection to mongoDB.",
			"error":  err,
		}).Info("Mongo repository info.")
		return err
	}
	log.WithFields(log.Fields{
		"time":   time.Now(),
		"method": "CreateUser()",
		"status": "Successfully connected to mongoDB.",
	}).Info("Mongo repository info.")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := client.Database("crudserver").Collection("users")

	result, err := col.InsertOne(ctx, bson.D{
		{Key: "_id", Value: uuid.New().String()},
		{Key: "userName", Value: u.UserName},
		{Key: "userAge", Value: u.UserAge},
		{Key: "isAdult", Value: u.IsAdult},
	})
	if err != nil {
		log.WithFields(log.Fields{
			"method": "CreateUser()",
			"status": "Failed while inserting.",
			"error":  err,
		}).Info("Mongo repository info.")
		return err
	}
	log.WithFields(log.Fields{
		"method":     "CreateUser()",
		"status":     "Successfully inserted.",
		"insertedID": result.InsertedID,
	}).Info("Mongo repository info.")
	return nil
}

// ReadUser returns User object from mongo database
// with selection by UserId
func (rps MongoRepository) ReadUser(userID string) (User, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.WithFields(log.Fields{
			"method": "ReadUser()",
			"status": "Failed connection to mongoDB.",
			"error":  err,
		}).Info("Mongo repository info.")
		return User{}, err
	}
	log.WithFields(log.Fields{
		"time":   time.Now(),
		"method": "ReadUser()",
		"status": "Successfully connected to mongoDB.",
	}).Info("Mongo repository info.")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := client.Database("crudserver").Collection("users")

	var rUser User
	err = col.FindOne(ctx, bson.D{{Key: "_id", Value: userID}}).Decode(&rUser)
	if err != nil {
		log.WithFields(log.Fields{
			"method": "ReadUser()",
			"status": "Failed while reading.",
			"error":  err,
		}).Info("Mongo repository info.")
		return User{}, err
	}
	log.WithFields(log.Fields{
		"method": "ReadUser()",
		"status": "Successfully read.",
		"time":   time.Now(),
	}).Info("Mongo repository info.")
	return rUser, nil
}

// UpdateUser update User object from mongo database
// with selection by UserId
func (rps MongoRepository) UpdateUser(u User) error {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.WithFields(log.Fields{
			"method": "UpdateUser()",
			"status": "Failed connection to mongoDB.",
			"error":  err,
		}).Info("Mongo repository info.")
		return err
	}
	log.WithFields(log.Fields{
		"time":   time.Now(),
		"method": "UpdateUser()",
		"status": "Successfully connected to mongoDB.",
	}).Info("Mongo repository info.")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := client.Database("crudserver").Collection("users")
	fmt.Println(u.UserID)
	result, err := col.UpdateOne(ctx, bson.D{{Key: "_id", Value: u.UserID}}, bson.D{
		{Key: "userName", Value: u.UserName},
		{Key: "userAge", Value: u.UserAge},
		{Key: "isAdult", Value: u.IsAdult},
	})
	if err != nil {
		log.WithFields(log.Fields{
			"method": "UpdateUser()",
			"status": "Failed while updating.",
			"error":  err,
		}).Info("Mongo repository info.")
		return err
	}
	log.WithFields(log.Fields{
		"method":        "UpdateUser()",
		"status":        "Successfully updated.",
		"upsertedCount": result.UpsertedCount,
	}).Info("Mongo repository info.")
	return nil
}

// DeleteUser delete User object from mongo database
// with selection by UserId
func (rps MongoRepository) DeleteUser(userID string) error {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.WithFields(log.Fields{
			"method": "DeleteUser()",
			"status": "Failed connection to mongoDB.",
			"error":  err,
		}).Info("Mongo repository info.")
		return err
	}
	log.WithFields(log.Fields{
		"time":   time.Now(),
		"method": "DeleteUser()",
		"status": "Successfully connected to mongoDB.",
	}).Info("Mongo repository info.")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := client.Database("crudserver").Collection("users")

	result, err := col.DeleteOne(ctx, bson.D{{Key: "_id", Value: userID}})
	if err != nil {
		log.WithFields(log.Fields{
			"method": "DeleteUser()",
			"status": "Failed while deleting.",
			"error":  err,
		}).Info("Mongo repository info.")
		return err
	}
	log.WithFields(log.Fields{
		"method":       "DeleteUser()",
		"status":       "Successfully deleted.",
		"deletedCount": result.DeletedCount,
	}).Info("Mongo repository info.")
	return nil
}

// AddImage function
func (rps MongoRepository) AddImage() {

}

// GetImage function
func (rps MongoRepository) GetImage() {

}

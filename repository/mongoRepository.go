package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
)

type MongoRepository struct{}

func(rps MongoRepository) CreateUser(u User) error{

	clientOptions := options.Client().ApplyURI("mogodb://27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil{
		log.WithFields(log.Fields{
			"method" : "CreateUser",
			"error" : err,
		}).Info("Mongodb repository info.")
		return err
	}
	defer client.Disconnect(context.TODO())

	col := client.Database("crudserver").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	result, err := col.InsertOne(ctx, u)
	if err != nil{
		log.WithFields(log.Fields{
			"method" : "CreateUser",
			"error" : err,
		}).Info("Mongodb repository info.")
		return err
	}else{
		log.WithFields(log.Fields{
			"method" : "CreateUser",
			"insertedID" : result.InsertedID,
		}).Info("Mongodb repository info.")
	}
	return nil
}

func(rps MongoRepository) ReadUser(){

	
}

func(rps MongoRepository) UpdateUser(u User)error{
	clientOptions := options.Client().ApplyURI("mogodb://27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil{
		log.WithFields(log.Fields{
			"method" : "CreateUser",
			"error" : err,
		}).Info("Mongodb repository info.")
		return err
	}
	defer client.Disconnect(context.TODO())

	col := client.Database("crudserver").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	result, err := col.UpdateByID(ctx, u.UserId, u)
	if err != nil{
		log.WithFields(log.Fields{
			"method" : "UpdateUser",
			"error" : err,
		}).Info("Mongodb repository info.")
		return err
	}else{
		log.WithFields(log.Fields{
			"method" : "DeleteUser",
			"upsertedId" : result.UpsertedID,
		}).Info("Mongodb repository info.")
	}
	return nil
}

func(rps MongoRepository) DeleteUser(userID uuid.UUID)error{
	clientOptions := options.Client().ApplyURI("mogodb://27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil{
		log.WithFields(log.Fields{
			"method" : "DeleteUser",
			"error" : err,
		}).Info("Mongodb repository info.")
		return err
	}
	defer client.Disconnect(context.TODO())

	col := client.Database("crudserver").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	result, err := col.DeleteOne(ctx, bson.M{"_id" : userID})
	if err != nil{
		log.WithFields(log.Fields{
			"method" : "DeleteUser",
			"error" : err,
		}).Info("Mongodb repository info.")
		return err
	}else{
		log.WithFields(log.Fields{
			"method" : "DeleteUser",
			"deleted count" : result.DeletedCount,
		}).Info("Mongodb repository info.")
	}

	return nil
}

func(rps MongoRepository) AddImage(){

}

func(rps MongoRepository) GetImage(){

}
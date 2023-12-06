package database


import (
	"context"
	"fmt"
	"log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Db *mongo.Database
var MongoClient *mongo.Client
var err error


func InitDB() (*mongo.Client, error) {
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	dbConnectionString := "mongodb+srv://hakulakaka:71cO1M0CC7VyVEFu@solitarygormetdb.r1ldxfd.mongodb.net/?retryWrites=true&w=majority"

	opts := options.Client().ApplyURI(dbConnectionString).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	MongoClient, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
		return nil, err

	}

	// Send a ping to confirm a successful connection
	err = MongoClient.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	dbName := "discordFileStorageServer"
	Db = MongoClient.Database(dbName)

	// Check if the database is not nil
	if Db == nil {
		return nil, fmt.Errorf("failed to find database %s", dbName)
	}

	fmt.Println("Connected to MongoDB!")

	err = UniqueIndexes()
	if err != nil {
		log.Fatal(err)
		fmt.Println("failed to set up index models")
		return nil, err
	}

	fmt.Println("Successfully setted up index models to avoid duplicates during insertion")

	return MongoClient, nil
}


func UniqueIndexes() error {
    collection1 := Db.Collection("fileList")

    indexmodel := mongo.IndexModel{
        Keys:    bson.D{{"name", 1}},
        Options: options.Index().SetUnique(true),
    }

    _, err := collection1.Indexes().CreateOne(context.TODO(), indexmodel)
    if err != nil {
        log.Fatal("Failed to create unique index:", err)
		return err
    }

	return nil
}

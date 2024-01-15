package database

import (
	"context"
	"fmt"
	"log"
	"github.com/hakuuww/hermione/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
        Keys:    bson.D{
						{"fileName", 1}, 
						{"fileType", 1},
					},
        Options: options.Index().SetUnique(true),
    }

    _, err := collection1.Indexes().CreateOne(context.TODO(), indexmodel)
    if err != nil {
        log.Fatal("Failed to create unique index:", err)
		return err
    }

	return nil
}

// InsertDocument inserts a document into the specified MongoDB collection
func InsertEmptyFileChunksDocument(ctx context.Context, collection *mongo.Collection, doc models.Document) error {
	_, err := collection.InsertOne(ctx, doc)
	return err
}

// AddFileChunkToDocument finds a document by ID and adds a new entry to the fileChunks array
func AddFileChunkToDocument(ctx context.Context, collection *mongo.Collection, docID primitive.ObjectID, fileChunk models.FileChunk) error {
	// Update the document with the new fileChunk
	filter := bson.M{"_id": docID}
	update := bson.M{"$push": bson.M{"fileChunks": fileChunk}}

	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

// GetFileChunksByFilename retrieves file chunks based on the provided filename.
func GetDocumentByFilename(ctx context.Context, fileList *mongo.Collection, filename string) (*models.Document, error) {
	// Search for the document in the collection based on the filename
	filter := bson.M{"fileName": filename}
	result := new(models.Document)
	err := fileList.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("error retrieving file document from db: %v", err)
	}

	return result, nil
}
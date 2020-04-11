package mongo

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(DBName, dbHost, dbPort string) *mongo.Client {
	host := "mongodb://" + dbHost + ":" + dbPort
	// Set client options
	clientOptions := options.Client().ApplyURI(host)

	// Connect to MongoDB
	dbClient, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = dbClient.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")
	return dbClient
}

func Find(collection *mongo.Collection, filter *bson.M) *mongo.Cursor {
	findOptions := options.Find()
	cur, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	// log.Println("Found documents ", cur)
	return cur
}

func FindOne(collection *mongo.Collection, filter *bson.M) *mongo.SingleResult {
	return collection.FindOne(context.TODO(), filter)
}

func Update(collection *mongo.Collection, filter *bson.D, update *bson.D) {
	updateResult, err := collection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
}

func Insert(collection *mongo.Collection, doc interface{}) {
	insertResult, err := collection.InsertOne(context.TODO(), doc)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Inserted a single document: ", insertResult.InsertedID)
}

func InsertMany(collection *mongo.Collection, docs []interface{}) {
	insertManyResult, err := collection.InsertMany(context.TODO(), docs)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)
}

func Close(dbClient *mongo.Client) {
	// Close the connection once no longer needed
	err := dbClient.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Connection to MongoDB closed.")
	}
}

func GetCollection(client *mongo.Client, DBName, collectionName string) *mongo.Collection {
	return client.Database(DBName).Collection(collectionName)
}

func ConnectCollection(db string, collectionName string) (*mongo.Client, *mongo.Collection) {
	mongoHost := "localhost"
	if mh := os.Getenv("MONGO_HOST"); mh != "" {
		mongoHost = mh
	}
	client := Connect(db, mongoHost, "27017")
	collection := GetCollection(client, db, collectionName)

	return client, collection
}

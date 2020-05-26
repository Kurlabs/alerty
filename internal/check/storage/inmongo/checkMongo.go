package inmongo

import (
	"context"
	"log"
	"sync"

	"github.com/Kurlabs/alerty/internal/check"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type monitorsRepository struct {
	monitors   map[primitive.ObjectID]check.Monitor
	collection *mongo.Collection
}

var (
	monitorsOnce     sync.Once
	monitorsInstance *monitorsRepository
)

func NewMonitorsRepository(collection *mongo.Collection) check.Repository {
	monitorsOnce.Do(func() {
		monitorsInstance = &monitorsRepository{
			monitors:   make(map[primitive.ObjectID]check.Monitor),
			collection: collection,
		}
	})
	return monitorsInstance
}

func (m *monitorsRepository) GetByID(ID string) (*check.Monitor, error) {
	var monitor check.Monitor

	monitorID, _ := primitive.ObjectIDFromHex(ID)
	filter := bson.M{
		"_id": monitorID,
	}
	opts := options.FindOne()
	err := m.collection.FindOne(context.TODO(), filter, opts).Decode(&monitor)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &monitor, nil
}

func (m *monitorsRepository) GetOne(filter *bson.M) (*check.Monitor, error) {
	var monitor check.Monitor
	opts := options.FindOne()
	err := m.collection.FindOne(context.TODO(), filter, opts).Decode(&monitor)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &monitor, nil
}

func (m *monitorsRepository) Find(filter *bson.M) (*mongo.Cursor, error) {
	findOptions := options.Find()
	cur, err := m.collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	return cur, nil
}

func (m *monitorsRepository) Delete(ID string) error {
	monitorID, _ := primitive.ObjectIDFromHex(ID)
	filter := bson.M{
		"_id": monitorID,
	}
	opts := options.Delete()
	_, err := m.collection.DeleteOne(context.TODO(), filter, opts)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (m *monitorsRepository) Save(monitor check.Monitor) error {
	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"_id", monitor.ID}}
	// data, err := json.Marshal(monitor)
	// fmt.Println(data)
	primitive.NewObjectID()

	update := bson.D{{"$set", monitor}}
	result, err := m.collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		log.Fatal(err)
		return err
	}

	if result.MatchedCount != 0 {
		log.Println("matched and replaced an existing document")
	}
	if result.UpsertedCount != 0 {
		log.Printf("inserted a new document with ID %v\n", result.UpsertedID)
	}

	return nil
}

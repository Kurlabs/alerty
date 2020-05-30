package mongo

import (
	"github.com/Kurlabs/alerty/shared/env"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	monitorsCOLLECTION     = "monitors"
	eventsCOLLECTION       = "events"
	messagesCOLLECTION     = "messages"
	contactsCOLLECTION     = "contacts"
	integrationsCOLLECTION = "integrations"
	metricsCOLLECTION      = "metrics"
)

var (
	client *mongo.Client
)

func init() {
	client = Connect(env.Config.MongoHost, env.Config.MongoPort)
}

// MCollection Monitor's collection
func MCollection() *mongo.Collection {
	return GetCollection(client, env.Config.DBName, monitorsCOLLECTION)
}

// ECollection Event's collection
func ECollection() *mongo.Collection {
	return GetCollection(client, env.Config.DBName, eventsCOLLECTION)
}

// MSCollection Messages' collection
func MSCollection() *mongo.Collection {
	return GetCollection(client, env.Config.DBName, messagesCOLLECTION)
}

// CCollection Contacts' collection
func CCollection() *mongo.Collection {
	return GetCollection(client, env.Config.DBName, contactsCOLLECTION)
}

// ICollection Integration's collection
func ICollection() *mongo.Collection {
	return GetCollection(client, env.Config.DBName, integrationsCOLLECTION)
}

// MTCollection Metrics' collection
func MTCollection() *mongo.Collection {
	return GetCollection(client, env.Config.DBName, metricsCOLLECTION)
}

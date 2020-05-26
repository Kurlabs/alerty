package inmongo

import (
	"context"
	"os"
	"testing"

	"github.com/Kurlabs/alerty/internal/check"
	models "github.com/Kurlabs/alerty/shared/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

func TestNewMonitor(t *testing.T) {
	DBName := "alerty"
	monitorsCollection := "monitors"
	mongoHost := "localhost"
	if mh := os.Getenv("MONGO_HOST"); mh != "" {
		mongoHost = mh
	}
	client := models.Connect(DBName, mongoHost, "27017")
	collection := models.GetCollection(client, DBName, monitorsCollection)
	checkRepo := NewMonitorsRepository(collection)

	website, err := check.New("Alerty", "https://alerty.online", 10, 1)
	if err != nil {
		t.Errorf("%v", err)
	}

	err = checkRepo.Save(*website)
	if err != nil {
		t.Errorf("%v", err)
	}

	website2, err := checkRepo.Get(website.ID.Hex())
	if err != nil {
		t.Logf("%v", website2)
		t.Errorf("%v", err)
	}

	website2, err = checkRepo.GetOne(&bson.M{"url": "https://alerty.online"})
	if err != nil {
		t.Logf("%v", website2)
		t.Errorf("%v", err)
	}

	err = checkRepo.Delete(website.ID.Hex())
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestFindAll(t *testing.T) {
	DBName := "alerty"
	monitorsCollection := "monitors"
	mongoHost := "localhost"
	if mh := os.Getenv("MONGO_HOST"); mh != "" {
		mongoHost = mh
	}
	client := models.Connect(DBName, mongoHost, "27017")
	collection := models.GetCollection(client, DBName, monitorsCollection)
	checkRepo := NewMonitorsRepository(collection)

	cur, err := checkRepo.Find(&bson.M{})
	if err != nil {
		t.Errorf("%v", err)
	}

	for cur.Next(context.TODO()) {
		var monitor check.Monitor
		err := cur.Decode(&monitor)
		if err != nil {
			t.Errorf("%v", err)
		}
		t.Logf("%v", monitor)
	}

}

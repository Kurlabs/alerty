package inmongo

import (
	"context"
	"testing"

	"github.com/Kurlabs/alerty/internal/check"
	models "github.com/Kurlabs/alerty/shared/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

func TestNewMonitor(t *testing.T) {
	collection := models.MCollection()
	checkRepo := NewMonitorsRepository(collection)

	website, err := check.New("Alerty", "https://alerty.online", 10, 1)
	if err != nil {
		t.Errorf("%v", err)
	}

	err = checkRepo.Save(*website)
	if err != nil {
		t.Errorf("%v", err)
	}

	website2, err := checkRepo.GetByID(website.ID.Hex())
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
	collection := models.MCollection()
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

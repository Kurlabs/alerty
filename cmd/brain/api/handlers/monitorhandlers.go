package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	conn "github.com/Kurlabs/alerty/shared/mongo"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var collection = conn.MCollection()

// RunnerStatus Message sent from runner monitor
type RunnerStatus struct {
	ID           primitive.ObjectID `json:"id"`
	URL          string             `json:"url"`
	Response     int16              `json:"response"`
	LastResponse int16              `json:"last_response"`
	RequestTime  float32            `json:"request_time"`
}

// RobotResponse robot resuts
type RobotResponse struct {
	ID          primitive.ObjectID `json:"id"`
	URL         string             `json:"url"`
	RobotResult string             `json:"robot_result"`
}

func MonitorBatch(c echo.Context) error {
	defer c.Request().Body.Close()
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("Falied reading the request body: %s", err)
		return c.String(http.StatusInternalServerError, "")
	}
	// Unmarshal
	var msgs []RunnerStatus
	err = json.Unmarshal(b, &msgs)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error")
	}

	// Update monitors collection
	updateBatch(msgs)

	// Marshal msgs to create a reponse
	output, err := json.Marshal(msgs)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error")
	}

	return c.String(http.StatusOK, string(output))
}

func MonitorRobot(c echo.Context) error {
	defer c.Request().Body.Close()
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error")
	}

	// Unmarshal
	var msgs []RobotResponse
	err = json.Unmarshal(b, &msgs)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error")
	}

	// Update monitors collection
	updateBatchRobot(msgs)

	// Marshal msgs to create a reponse
	output, err := json.Marshal(msgs)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error")
	}

	return c.String(http.StatusOK, string(output))

}

// updateBatch recieve a array of runners and update monitor collection
func updateBatch(msgs []RunnerStatus) {
	var controlled bool
	for _, msg := range msgs {
		filter := bson.D{{"_id", msg.ID}}
		log.Printf("URL: %s ", msg.URL)
		if msg.Response != msg.LastResponse {
			controlled = false
			update := bson.D{
				{"$set", bson.D{
					{"response", msg.Response},
					{"last_response", msg.LastResponse},
					{"controlled", controlled},
					{"request_time", msg.RequestTime},
				}},
				{"$currentDate", bson.D{
					{"updated_at", true},
				}},
			}
			conn.Update(collection, &filter, &update)
		} else {
			update := bson.D{
				{"$set", bson.D{
					{"response", msg.Response},
					{"last_response", msg.LastResponse},
					{"request_time", msg.RequestTime},
				}},
				{"$currentDate", bson.D{
					{"updated_at", true},
				}},
			}
			conn.Update(collection, &filter, &update)
		}
	}
}

// updateBatch recieve a array of runners and update monitor collection
func updateBatchRobot(msgs []RobotResponse) {
	for _, msg := range msgs {
		filter := bson.D{{"_id", msg.ID}}
		log.Printf("ID: %s ", msg.ID)
		update := bson.D{
			{"$set", bson.D{
				{"robot_result", msg.RobotResult},
			}},
		}
		log.Printf("URL: %s ", msg.URL)
		conn.Update(collection, &filter, &update)
	}
}

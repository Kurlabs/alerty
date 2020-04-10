package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/Kurlabs/alerty/shared/env"
	conn "github.com/Kurlabs/alerty/shared/mongo"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var collection *mongo.Collection

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

func main() {
	mongoHost := "localhost"
	if mh := os.Getenv("MONGO_HOST"); mh != "" {
		mongoHost = mh
	}
	dbclient := conn.Connect(env.Config.DBName, mongoHost, env.Config.MongoPort)
	collection = conn.GetCollection(dbclient, env.Config.DBName, env.Config.MonitorCollection)

	// Define routers
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/"))
	})

	r.Route("/monitors/", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("/"))
		})
		r.Post("/batch", Batch)
		r.Post("/robot", Robot)
	})

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		fmt.Printf("%s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		fmt.Printf("Logging err: %s\n", err.Error())
	}

	http.ListenAndServe(":3000", r)
}

// updateBatch recieve a array of runners and update monitor collection
func updateBatch(msgs []RunnerStatus) {
	var controlled bool
	for _, msg := range msgs {
		filter := bson.D{{"_id", msg.ID}}
		fmt.Printf("URL: %s ", msg.URL)
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

// Batch returns done status
func Batch(w http.ResponseWriter, r *http.Request) {
	// Validate token
	token := r.Header.Get("TOKEN")
	if token != env.Config.BrainToken {
		http.Error(w, token, 403)
		return
	}
	// Read body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Unmarshal
	var msgs []RunnerStatus
	err = json.Unmarshal(b, &msgs)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Update monitors collection
	updateBatch(msgs)

	// Marshal msgs to create a reponse
	output, err := json.Marshal(msgs)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(output)
}

// updateBatch recieve a array of runners and update monitor collection
func updateBatchRobot(msgs []RobotResponse) {
	for _, msg := range msgs {
		filter := bson.D{{"_id", msg.ID}}
		fmt.Printf("ID: %s ", msg.ID)
		update := bson.D{
			{"$set", bson.D{
				{"robot_result", msg.RobotResult},
			}},
		}
		fmt.Printf("URL: %s ", msg.URL)
		conn.Update(collection, &filter, &update)
	}
}

// Robot returns done status
func Robot(w http.ResponseWriter, r *http.Request) {
	// Validate token
	token := r.Header.Get("TOKEN")
	if token != env.Config.BrainToken {
		http.Error(w, token, 403)
		return
	}
	// Read body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Unmarshal
	var msgs []RobotResponse
	err = json.Unmarshal(b, &msgs)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Update monitors collection
	updateBatchRobot(msgs)

	// Marshal msgs to create a reponse
	output, err := json.Marshal(msgs)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(output)
}

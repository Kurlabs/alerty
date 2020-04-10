package check

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Monitor struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CLS        string             `json:"cls" bson:"_cls"`
	User       primitive.ObjectID
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
	Timeout    int
	Name       string
	Status     bool
	Validated  bool
	Frequency  int
	SSL        bool
	URL        string
	Response   int
	Updated    time.Time
	Controlled bool
	IP         string
	Port       int
	Path       string
	Protocol   string
}

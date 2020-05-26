package check

import (
	"log"
	"net/url"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Monitor struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CLS       string             `json:"cls" bson:"_cls"` // Monitor.WebsiteMonitor or Monitor.SocketMonitor
	User      primitive.ObjectID
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	Updated   time.Time
	// Website
	URL string
	// Socket
	IP       string
	Port     int
	Path     string
	Protocol string // tpc or udp
	// Common
	SSL       bool
	Name      string
	Timeout   int
	Frequency int
	// Logic
	Controlled bool
	Status     bool
	Validated  bool
	Response   int
}

const (
	clsWebsite       = "Monitor.WebsiteMonitor"
	clsSocket        = "Monitor.SocketMonitor"
	controlDefault   = false
	statusDefault    = false
	validatedDefault = false
)

// New creates a Website or Socket Monitor with the minimun data
func New(name, URL string, timeout, frequency int) (*Monitor, error) {
	var SSL bool
	var monitorType string
	var port int

	u, err := url.Parse(URL)
	if err != nil {
		log.Fatalf("%v", err)
		return nil, err
	}

	if u.Scheme == "http" {
		SSL = false
		port = 80
		monitorType = clsWebsite
	} else if u.Scheme == "https" {
		SSL = true
		port = 443
		monitorType = clsWebsite
	} else {
		// TODO : support SSL in sockets
		SSL = false
		port, err = strconv.Atoi(u.Port())
		if err != nil {
			return nil, err
		}
		monitorType = clsSocket
	}

	w := Monitor{
		ID:        primitive.NewObjectID(),
		CLS:       monitorType,
		User:      primitive.NilObjectID,
		CreatedAt: time.Now(),
		Updated:   time.Now(),
		// Website
		URL: URL,
		// Socket
		IP:       u.Hostname(),
		Port:     port,
		Path:     u.RawPath,
		Protocol: u.Scheme,
		// Common
		Name:      name,
		SSL:       SSL,
		Timeout:   timeout,
		Frequency: frequency,
		// Logic
		Controlled: controlDefault,
		Status:     statusDefault,
		Validated:  validatedDefault,
		Response:   0,
	}

	return &w, nil
}

type Repository interface {
	GetByID(ID string) (*Monitor, error)
	GetOne(filter *bson.M) (*Monitor, error)
	Save(monitor Monitor) error
	Find(filter *bson.M) (*mongo.Cursor, error)
	Delete(ID string) error
}

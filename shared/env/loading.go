package env

import (
	"log"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

type config struct {
	DBName            string
	MonitorCollection string
	EventCollection   string
	Level             string
	BrainURL          string
	BrainToken        string
	MongoPort         string
	MongoHost         string
}

// Config contains global variables
var (
	Config config
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// This is necesarry to run test
// https://github.com/joho/godotenv/issues/43#issuecomment-503183127
const projectDirName = "alerty"

func getEnvPath() string {
	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))
	return string(rootPath) + `/.env`
}

// use godot package to load/read the .env file and
func init() {
	log.Println("Loading .env file...")
	envPath := getEnvPath()
	err := godotenv.Load(envPath)

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	if os.Getenv("project_id") == "" {
		log.Fatalf("project_id was expected in .env file")
	}
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		log.Fatalf("GOOGLE_APPLICATION_CREDENTIALS was expected in .env file")
	}
	if os.Getenv("DBName") == "" {
		log.Fatalf("DBName was expected in .env file")
	}
	if os.Getenv("BrainURL") == "" {
		log.Fatalf("BrainURL was expected in .env file")
	}
	if os.Getenv("BrainToken") == "" {
		log.Fatalf("BrainToken was expected in .env file")
	}
	Config = config{
		DBName:            getEnv("DBName", "alerty"),
		MonitorCollection: "monitors",
		EventCollection:   "events",
		Level:             getEnv("Level", "debug"),
		BrainURL:          os.Getenv("BrainURL"),
		BrainToken:        os.Getenv("BrainToken"),
		MongoPort:         getEnv("MongoPort", "27017"),
		MongoHost:         getEnv("MongoHost", "localhost"),
	}
}

func (c *config) getDBName() string {
	return c.DBName
}

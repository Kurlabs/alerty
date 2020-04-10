package check

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	APP_TYPE   = "app"
	EMAIL_TYPE = "email"
	SMS_TYPE   = "sms"
	SLACK_TYPE = "slack"

	INFO_LEVEL     = "info"
	WARNING_LEVEL  = "warning"
	CRITICAL_LEVEL = "critical"
	DONE_LEVEL     = "done"
)

type Metric struct {
	Name string
}

type Rule struct {
	Metric   primitive.ObjectID
	Operator string
	Value    string
}

type Contact struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CLS           string             `json:"cls" bson:"_cls"`
	Name          string
	LastName      string `json:"last_name" bson:"last_name"`
	Email         string
	Area          int
	Number        int
	ContactParent primitive.ObjectID `json:"contact_parent" bson:"contact_parent"`
}

type Integration struct {
	ID                  primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CLS                 string             `json:"cls" bson:"_cls"`
	AccessToken         string             `json:"access_token" bson:"access_token"`
	Scope               string             `json:"scope" bson:"scope"`
	SlackUserID         string             `json:"slack_user_id" bson:"slack_user_id"`
	TeamName            string             `json:"team_name" bson:"team_name"`
	TeamID              string             `json:"team_id" bson:"team_id"`
	EnterpriseID        string             `json:"enterprise_id" bson:"enterprise_id"`
	IncWebhookChannel   string             `json:"inc_webhook_channel" bson:"inc_webhook_channel"`
	IncWebhookChannelID string             `json:"inc_webhook_channel_id" bson:"inc_webhook_channel_id"`
	IncWebhookConfigURL string             `json:"inc_webhook_config_url" bson:"inc_webhook_config_url"`
	IncWebhookURL       string             `json:"inc_webhook_url" bson:"inc_webhook_url"`
	BotUserID           string             `json:"bot_user_id" bson:"bot_user_id"`
	BotAccessToken      string             `json:"bot_access_token" bson:"bot_access_token"`
	SlackResponse       string             `json:"slack_response" bson:"slack_response"`
}

type Event struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Monitor      primitive.ObjectID
	Rules        []Rule
	Contacts     []primitive.ObjectID
	Integrations []primitive.ObjectID
}

type Message struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Sent            bool
	MessageType     string
	Level           string
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
	User            primitive.ObjectID
	Event           primitive.ObjectID
	Contact         primitive.ObjectID
	MonitorName     string
	MonitorURL      string
	MonitorPath     string
	MonitorCLS      string
	MonitorID       string
	ContactName     string
	IntegrationName string
}

type EventHistoryEntry struct {
	ID      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Event   primitive.ObjectID
	Monitor primitive.ObjectID
}

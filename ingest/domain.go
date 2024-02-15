package ingest

import (
	"encoding/json"
	"github.com/google/uuid"
	rabbitMQ "github.com/julioisaac/users/providers/rabbitmq"
)

const (
	appName         = "user_ingest"
	usersQueue      = "users.events.q"
	usersExchange   = "users.events.x"
	usersRoutingKey = "users.events.k"
)

type UserField int

const (
	ID UserField = iota
	FirstName
	LastName
	EmailAddress
	CreatedAt
	DeletedAt
	MergedAt
	ParentUserID
)

type User struct {
	ID           string `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	EmailAddress string `json:"email_address"`
	CreatedAt    string `json:"created_at"`
	DeletedAt    string `json:"deleted_at"`
	MergedAt     string `json:"merged_at"`
	ParentUserID string `json:"parent_user_id"`
}

func (u *User) ToRabbitEvent(app string, queue string) (rabbitMQ.EventDTO, error) {
	rabbitEvent := rabbitMQ.EventDTO{}

	uuid := uuid.NewString()

	rabbitEvent.MessageID = uuid
	rabbitEvent.Operation = uuid

	body, err := json.Marshal(u)

	if err != nil {
		return rabbitEvent, err
	}

	rabbitEvent.Body = body

	rabbitEvent.App = app
	rabbitEvent.Queue = queue

	return rabbitEvent, nil
}

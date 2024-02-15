package ingest

import (
	"encoding/csv"
	"fmt"
	"github.com/julioisaac/users/config"
	"github.com/julioisaac/users/logger"
	rabbitMQ "github.com/julioisaac/users/providers/rabbitmq"
	"os"
)

// @title Ingest users
// @version 1.0
// @description Ingest users application

type Service interface {
	Run(filePath string) error
}

type ingestService struct {
	rabbitProvider rabbitMQ.RabbitMQProvider
}

func NewService(rabbitProvider rabbitMQ.RabbitMQProvider) Service {
	return &ingestService{
		rabbitProvider: rabbitProvider,
	}
}

// Run ingestion of users
func (is *ingestService) Run(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening the file:", err)
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	lines, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV file:", err)
		return err
	}

	var events []*rabbitMQ.EventDTO
	for i, line := range lines {
		if i == 0 {
			continue
		}

		usr := buildUser(line)
		event, err := usr.ToRabbitEvent(appName, usersQueue)
		if err != nil {
			logger.Logger.Errorf("[Ingest][sendEvent][ERROR]: %s", err.Error())
			return err
		}
		events = append(events, &event)

		if len(events) == config.GetInt("BATCH_SIZE") {
			is.sendEventBatch(events)
			events = nil
		}
	}
	if len(events) > 0 {
		is.sendEventBatch(events)
	}
	return nil
}

func (is *ingestService) sendEventBatch(events []*rabbitMQ.EventDTO) {
	is.rabbitProvider.SetupProducer(rabbitMQ.Producer{
		Name:       appName,
		Exchange:   usersExchange,
		RoutingKey: usersRoutingKey,
	})
	is.rabbitProvider.Producer(events)
}

func buildUser(line []string) User {
	return User{
		ID:           line[ID],
		FirstName:    line[FirstName],
		LastName:     line[LastName],
		EmailAddress: line[EmailAddress],
		CreatedAt:    line[CreatedAt],
		DeletedAt:    line[DeletedAt],
		MergedAt:     line[MergedAt],
		ParentUserID: line[ParentUserID],
	}
}

//go:build unit

package ingest

import (
	"encoding/csv"
	"github.com/google/uuid"
	rabbitMQ "github.com/julioisaac/users/providers/rabbitmq"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type MockRabbitMQProvider struct {
	ProducerCalls [][]*rabbitMQ.EventDTO
}

func (m *MockRabbitMQProvider) SetupProducer(producer rabbitMQ.Producer) rabbitMQ.RabbitMQProvider {
	return nil
}

func (m *MockRabbitMQProvider) SetupConsumer(consumer rabbitMQ.Consumer) rabbitMQ.RabbitMQProvider {
	return nil
}

func (m *MockRabbitMQProvider) Producer(events []*rabbitMQ.EventDTO) {
	m.ProducerCalls = append(m.ProducerCalls, events)
}

func (m *MockRabbitMQProvider) Consumer() error {
	return nil
}

func TestRun(t *testing.T) {
	t.Run("BatchEvents", func(t *testing.T) {
		mockRabbitMQ := &MockRabbitMQProvider{}
		ingestService := NewService(mockRabbitMQ)

		tempFile, err := os.CreateTemp("", "test_batch_events_*.csv")
		assert.NoError(t, err)
		defer os.Remove(tempFile.Name())

		csvWriter := csv.NewWriter(tempFile)
		defer csvWriter.Flush()

		header := []string{"id", "first_name", "last_name", "email_address", "created_at", "deleted_at", "merged_at", "parent_user_id"}
		assert.NoError(t, csvWriter.Write(header))

		for i := 0; i < 10; i++ {
			line := []string{uuid.NewString(), "John", "Doe", "john.doe@example.com", "1609459200000", "-1", "-1", "-1"}
			assert.NoError(t, csvWriter.Write(line))
		}

		csvWriter.Flush()
		tempFile.Close()

		assert.NoError(t, ingestService.Run(tempFile.Name()))

		assert.Len(t, mockRabbitMQ.ProducerCalls, 1)
		assert.Len(t, mockRabbitMQ.ProducerCalls[0], 10)
	})

	t.Run("SingleEvent", func(t *testing.T) {
		mockRabbitMQ := &MockRabbitMQProvider{}
		ingestService := NewService(mockRabbitMQ)

		tempFile, err := os.CreateTemp("", "test_single_event_*.csv")
		assert.NoError(t, err)
		defer os.Remove(tempFile.Name())

		csvWriter := csv.NewWriter(tempFile)
		defer csvWriter.Flush()

		header := []string{"id", "first_name", "last_name", "email_address", "created_at", "deleted_at", "merged_at", "parent_user_id"}
		assert.NoError(t, csvWriter.Write(header))

		line := []string{uuid.NewString(), "Jane", "Doe", "jane.doe@example.com", "1609459200000", "-1", "-1", "-1"}
		assert.NoError(t, csvWriter.Write(line))

		csvWriter.Flush()
		tempFile.Close()

		assert.NoError(t, ingestService.Run(tempFile.Name()))

		assert.Len(t, mockRabbitMQ.ProducerCalls, 1)
		assert.Len(t, mockRabbitMQ.ProducerCalls[0], 1)
	})

	t.Run("NoEvents", func(t *testing.T) {
		mockRabbitMQ := &MockRabbitMQProvider{}
		ingestService := NewService(mockRabbitMQ)

		tempFile, err := os.CreateTemp("", "test_no_events_*.csv")
		assert.NoError(t, err)
		defer os.Remove(tempFile.Name())

		assert.NoError(t, ingestService.Run(tempFile.Name()))

		assert.Len(t, mockRabbitMQ.ProducerCalls, 0)
	})
}

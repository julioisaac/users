//go:build unit

package ingest

import (
	"testing"

	rabbitMQ "github.com/julioisaac/users/providers/rabbitmq"
)

func TestToRabbitEvent(t *testing.T) {
	tests := []struct {
		name     string
		userData *User
		expected rabbitMQ.EventDTO
	}{
		{
			name: "ValidUser",
			userData: &User{
				ID:           "123",
				FirstName:    "John",
				LastName:     "Doe",
				EmailAddress: "john.doe@example.com",
				CreatedAt:    "1609459200000",
				DeletedAt:    "-1",
				MergedAt:     "-1",
				ParentUserID: "-1",
			},
			expected: rabbitMQ.EventDTO{
				Operation: "ingest",
				Body:      []byte(`{"id":"123","first_name":"John","last_name":"Doe","email_address":"john.doe@example.com","created_at":"1609459200000","deleted_at":"-1","merged_at":"-1","parent_user_id":"-1"}`),
				App:       appName,
				Queue:     usersQueue,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			result, err := tt.userData.ToRabbitEvent(appName, usersQueue)

			if err != nil {
				t.Fatalf("Unexpected error for test %s: %v", tt.name, err)
			}

			if string(result.Body) != string(tt.expected.Body) ||
				result.App != tt.expected.App ||
				result.Queue != tt.expected.Queue {
				t.Errorf("Test %s failed. \n Expected: %+v, \n Got: %+v", tt.name, tt.expected, result)
			}
		})
	}
}

//go:build unit

package consumer

import (
	"github.com/julioisaac/users/user"
	"testing"
	"time"
)

func TestToEntity(t *testing.T) {
	tests := []struct {
		name     string
		userData *User
		expected user.User
	}{
		{
			name: "ValidUser",
			userData: &User{
				ID:           "8",
				FirstName:    "Hanah",
				LastName:     "Schmidt",
				EmailAddress: "Hanah_Schmidt1965@gmail.edu",
				CreatedAt:    "1361218223000",
				DeletedAt:    "-1",
				MergedAt:     "-1",
				ParentUserID: "-1",
			},
			expected: user.User{
				UserID:       int64(8),
				FirstName:    "Hanah",
				LastName:     "Schmidt",
				EmailAddress: "Hanah_Schmidt1965@gmail.edu",
				CreatedAt:    time.Unix(0, 1361218223000*int64(time.Millisecond)),
				DeletedAt:    nil,
				MergedAt:     nil,
				ParentUserID: int64(-1),
			},
		},
		{
			name: "UserWithDeletedAt",
			userData: &User{
				ID:           "10",
				FirstName:    "John",
				LastName:     "Doe",
				EmailAddress: "john.doe@example.com",
				CreatedAt:    "1370012345000",
				DeletedAt:    "1380012345000",
				MergedAt:     "-1",
				ParentUserID: "-1",
			},
			expected: user.User{
				UserID:       int64(10),
				FirstName:    "John",
				LastName:     "Doe",
				EmailAddress: "john.doe@example.com",
				CreatedAt:    time.Unix(0, 1370012345000*int64(time.Millisecond)),
				DeletedAt: func() *time.Time {
					t := time.Unix(0, 1380012345000*int64(time.Millisecond))
					return &t
				}(),
				MergedAt:     nil,
				ParentUserID: int64(-1),
			},
		},
		{
			name: "UserWithMergedAt",
			userData: &User{
				ID:           "12",
				FirstName:    "Alice",
				LastName:     "Smith",
				EmailAddress: "alice.smith@example.com",
				CreatedAt:    "1380012345000",
				DeletedAt:    "-1",
				MergedAt:     "1390012345000",
				ParentUserID: "-1",
			},
			expected: user.User{
				UserID:       int64(12),
				FirstName:    "Alice",
				LastName:     "Smith",
				EmailAddress: "alice.smith@example.com",
				CreatedAt:    time.Unix(0, 1380012345000*int64(time.Millisecond)),
				DeletedAt:    nil,
				MergedAt: func() *time.Time {
					t := time.Unix(0, 1390012345000*int64(time.Millisecond))
					return &t
				}(),
				ParentUserID: int64(-1),
			},
		},
		{
			name: "UserWithParentUserID",
			userData: &User{
				ID:           "15",
				FirstName:    "Bob",
				LastName:     "Johnson",
				EmailAddress: "bob.johnson@example.com",
				CreatedAt:    "1400012345000",
				DeletedAt:    "-1",
				MergedAt:     "-1",
				ParentUserID: "8",
			},
			expected: user.User{
				UserID:       int64(15),
				FirstName:    "Bob",
				LastName:     "Johnson",
				EmailAddress: "bob.johnson@example.com",
				CreatedAt:    time.Unix(0, 1400012345000*int64(time.Millisecond)),
				DeletedAt:    nil,
				MergedAt:     nil,
				ParentUserID: int64(8),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.userData.ToEntity()

			if err != nil {
				t.Fatalf("Unexpected error for test %s: %v", tt.name, err)
			}

			if !timesAreEqual(result.DeletedAt, tt.expected.DeletedAt) || !timesAreEqual(result.MergedAt, tt.expected.MergedAt) {
				t.Errorf("Test %s failed. \n Expected: %+v, \n Got: %+v", tt.name, tt.expected, result)
				return
			}

			if result.UserID != tt.expected.UserID || result.ParentUserID != tt.expected.ParentUserID {
				t.Errorf("Test %s failed. \n Expected: %+v, \n Got: %+v", tt.name, tt.expected, result)
			}
		})
	}
}

func timesAreEqual(t1, t2 *time.Time) bool {
	if t1 == nil && t2 == nil {
		return true
	}
	if (t1 == nil && t2 != nil) || (t1 != nil && t2 == nil) {
		return false
	}
	return t1.Equal(*t2)
}

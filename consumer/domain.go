package consumer

import (
	"github.com/julioisaac/users/user"
	"strconv"
	"time"
)

const (
	appName    = "user_consumer"
	usersQueue = "users.events.q"
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

func (u *User) ToEntity() (user.User, error) {
	createdAt, err := parseTimestamp(u.CreatedAt)
	if err != nil {
		return user.User{}, err
	}

	deletedAt, err := parseOptionalTimestamp(u.DeletedAt)
	if err != nil {
		return user.User{}, err
	}

	mergedAt, err := parseOptionalTimestamp(u.MergedAt)
	if err != nil {
		return user.User{}, err
	}

	userID, err := strconv.ParseInt(u.ID, 10, 64)
	if err != nil {
		return user.User{}, err
	}

	parentUserID, err := strconv.ParseInt(u.ParentUserID, 10, 64)
	if err != nil {
		return user.User{}, err
	}

	return user.User{
		UserID:       userID,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		EmailAddress: u.EmailAddress,
		CreatedAt:    createdAt,
		DeletedAt:    deletedAt,
		MergedAt:     mergedAt,
		ParentUserID: parentUserID,
	}, nil
}

func parseTimestamp(timestamp string) (time.Time, error) {
	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(0, ts*int64(time.Millisecond)), nil
}

func parseOptionalTimestamp(timestamp string) (*time.Time, error) {
	if timestamp == "" || timestamp == "-1" {
		return nil, nil
	}

	ts, err := parseTimestamp(timestamp)
	if err != nil {
		return nil, err
	}

	return &ts, nil
}

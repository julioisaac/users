package user

import (
	"encoding/json"
	"time"
)

type User struct {
	Id           int64
	UserID       int64
	FirstName    string
	LastName     string
	EmailAddress string
	CreatedAt    time.Time
	DeletedAt    *time.Time
	MergedAt     *time.Time
	ParentUserID int64
}

type Users []User

func (u Users) ToBytes() ([]byte, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (u Users) First() User {
	return u[0]
}

// Pagination
type Pagination struct {
	Page  uint
	Limit uint
}

// Filters
type Filters struct {
	*Pagination  `json:"pagination"`
	UserID       int64  `json:"user_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	EmailAddress string `json:"email_address"`
	ParentUserID int64  `json:"parent_user_id"`
}

func (u Filters) ToBytes() ([]byte, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}
	return b, nil
}

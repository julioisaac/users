//go:build integration

package user_test

import (
	"context"
	"github.com/julioisaac/users/test"
	users "github.com/julioisaac/users/user"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var repo users.Repository

func TestMain(m *testing.M) {

	// Getting DBInstance
	dbInstance, _ := test.GetDBInstance()

	// Getting repository
	repo = users.NewUserRepo(dbInstance)

	// Running tests
	code := m.Run()

	os.Exit(code)

}

func Test_repository_Create_Find(t *testing.T) {

	t.Run("should create and find user", func(t *testing.T) {
		// given
		ctx := context.Background()
		user := users.User{
			UserID:       12,
			FirstName:    "Julio",
			LastName:     "Penha",
			EmailAddress: "julioisaac7@gmail.com",
			CreatedAt:    time.Now(),
		}

		// when
		err := repo.Create(ctx, user)
		if err != nil {
			t.Errorf("unexpected error creating user %d, err: %s", user.UserID, err)
			return
		}

		filters := users.Filters{
			Pagination: &users.Pagination{
				Page:  0,
				Limit: 10,
			},
			UserID: user.UserID,
		}

		u, err := repo.GetUserByParams(ctx, filters)
		if err != nil {
			t.Errorf("unexpected error finding user %d, err: %s", user.UserID, err)
			return
		}

		// then
		usr := u.First()
		assert.Equalf(t, int64(12), usr.UserID, "UserID")
		assert.Equalf(t, "Julio", usr.FirstName, "FirstName")
		assert.Equalf(t, "Penha", usr.LastName, "LastName")
		assert.Equalf(t, "julioisaac7@gmail.com", usr.EmailAddress, "EmailAddress")

	})

}

package ioc

import (
	users "github.com/julioisaac/users/user"
	"sync"
)

var (
	userRepoOnce sync.Once
	userRepo     users.Repository
)

// UserRepository with dependencies
func UserRepository() users.Repository {
	userRepoOnce.Do(func() {
		userRepo = users.NewUserRepo(DB())
	})

	return userRepo
}

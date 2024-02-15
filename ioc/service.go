package ioc

import (
	"github.com/julioisaac/users/ingest"
	"github.com/julioisaac/users/user"
	"sync"
)

var (
	ingestServiceOnce sync.Once
	ingestService     ingest.Service

	userServiceOnce sync.Once
	userService     user.Service
)

func IngestService() ingest.Service {
	ingestServiceOnce.Do(func() {
		ingestService = ingest.NewService(RabbitMQ())
	})

	return ingestService
}

func UserService() user.Service {
	userServiceOnce.Do(func() {
		userService = user.NewService(UserRepository(), RedisCache())
	})

	return userService
}

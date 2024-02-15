package user

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/julioisaac/users/logger"
	"github.com/julioisaac/users/providers/cache"
	"github.com/sirupsen/logrus"
)

type Service interface {
	Save(ctx context.Context, user User) error
	FindByFilters(ctx context.Context, filters Filters) (Users, error)
}

type service struct {
	repo  Repository
	cache cache.Cache
}

// NewService given dependencies
func NewService(repo Repository, cache cache.Cache) Service {
	return &service{
		repo:  repo,
		cache: cache,
	}
}

func (s service) Save(ctx context.Context, user User) error {

	err := s.repo.Create(ctx, user)
	if err != nil {
		logger.Logger.Errorf("cannot create user: %v", user)
		return err
	}

	return nil
}

func (s service) FindByFilters(ctx context.Context, filters Filters) (Users, error) {

	f, err := filters.ToBytes()
	if err != nil {
		logger.Logger.Warningf("error converting filters: %v, error: %v", f, err)
		return nil, err
	}

	filtersMD5 := md5.Sum(f)
	cacheKey := fmt.Sprintf("%x", filtersMD5)

	cachedUser, err := s.cache.Get(ctx, cacheKey)
	if !errors.Is(err, redis.Nil) && err != nil {
		logger.Logger.Warningf("error getting user by key: %s from cache, error: %v", cacheKey, err)
	}

	if len(cachedUser) != 0 {
		logger.Logger.Debugf("query: %v already exists in cache", filters)

		users := Users{}
		err = json.Unmarshal(cachedUser, &users)
		if err != nil {
			logger.Logger.Errorf("failed to convert bytes to Users struct, value: %v err: %s", cachedUser, err)
			return nil, err
		}
		return users, nil
	}

	users, err := s.repo.GetUserByParams(ctx, filters)
	if err != nil {
		logger.Logger.Errorf("error getting user by filters: %v err: %s", filters, err)
		return nil, err
	}

	u, err := users.ToBytes()
	if err != nil {
		logger.Logger.Warningf("error converting users: %v, error: %v", f, err)
		return nil, err
	}

	err = s.cache.Add(ctx, cacheKey, u, -1)
	if err != nil {
		logrus.Warningf("error adding query by filters: %v in cache", filters)
	} else {
		logrus.Debugf("query with filters: %v was added to cache", filters)
	}

	return users, nil
}

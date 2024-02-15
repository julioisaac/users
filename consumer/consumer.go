package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/julioisaac/users/ioc"
	"github.com/julioisaac/users/logger"
	rabbitMQ "github.com/julioisaac/users/providers/rabbitmq"
)

// @title Consumer users
// @version 1.0
// @description Consumer users application

// Start consumer of users
func Start() error {
	rp := rabbitMQ.NewRabbit()

	rp.SetupConsumer(rabbitMQ.Consumer{
		Queue:          usersQueue,
		Name:           appName,
		ProcessMessage: processMsg,
	})

	return rp.Consumer()
}

func processMsg(ctx context.Context, value []byte) error {
	user := User{}
	err := json.Unmarshal(value, &user)
	if err != nil {
		return fmt.Errorf("failed converting event to user, value: %s err: %s", string(value), err)
	}

	ue, err := user.ToEntity()
	if err != nil {
		return fmt.Errorf("failed converting user to user entity, user: %v err: %s", user, err)
	}

	userService := ioc.UserService()
	err = userService.Save(ctx, ue)
	if err != nil {
		logger.Logger.Errorf("Error saving user: %v", user)
		return err
	}

	logger.Logger.Infof("Processed message: %v", string(value))
	return nil
}

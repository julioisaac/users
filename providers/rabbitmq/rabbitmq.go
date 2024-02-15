package rabbitMQ

import (
	"context"
	"errors"
	"fmt"
	"github.com/julioisaac/users/logger"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/julioisaac/users/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQProvider interface {
	SetupProducer(pd Producer) RabbitMQProvider
	Producer(message []*EventDTO)
	SetupConsumer(c Consumer) RabbitMQProvider
	Consumer() error
}

type rabbitMQProvider struct {
	uri      string
	consumer Consumer
	producer Producer
}

// Consumer holds info to consume
type Consumer struct {
	Name           string
	Queue          string
	ProcessMessage func(ctx context.Context, value []byte) error
}

// Producer holds info to produce
type Producer struct {
	Name       string
	Exchange   string
	RoutingKey string
}

func NewRabbit() RabbitMQProvider {
	return &rabbitMQProvider{
		uri: config.GetString("RABBITMQ_URI"),
	}
}

func (p *rabbitMQProvider) SetupProducer(pd Producer) RabbitMQProvider {
	p.producer = pd
	return p
}

func (p *rabbitMQProvider) SetupConsumer(c Consumer) RabbitMQProvider {
	p.consumer = c
	return p
}

func (p *rabbitMQProvider) Producer(messages []*EventDTO) {

	conn, err := amqp.Dial(p.uri)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	msgCh := make(chan *EventDTO, len(messages))
	doneCh := make(chan struct{})

	workers := runtime.NumCPU() * 4

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go publishMessages(ctx, ch, &wg, p.producer, msgCh)
	}

	for _, msg := range messages {
		msgCh <- msg
	}

	go func() {
		wg.Wait()
		close(doneCh)
	}()

	<-doneCh
	close(msgCh)

	fmt.Println("Batch publishing complete")
}

func (p *rabbitMQProvider) Consumer() error {
	conn, err := amqp.Dial(p.uri)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	ctx := context.Background()

	msgs, err := ch.Consume(
		p.consumer.Queue,
		p.consumer.Name,
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-stopChan:
			return errors.New("received interrupt signal. Exiting")
		case msg := <-msgs:
			err = p.consumer.ProcessMessage(ctx, msg.Body)
			if err != nil {
				logger.Logger.Errorf("Failed consuming, err: %v", err)
			} else {
				err = msg.Ack(false)
				failOnError(err, "Failed to acknowledge the message")
			}
		}
	}

}

func publishMessages(ctx context.Context, ch *amqp.Channel, wg *sync.WaitGroup, pd Producer, msgCh chan *EventDTO) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-msgCh:
			if !ok {
				return
			}
			logger.Logger.Infof("[rabbbitmq][publishMessage][INFO] Publishing event in queue - %+v", string(msg.Body))

			err := publishMessage(ctx, ch, pd, msg)
			if err != nil {
				logger.Logger.Errorf("[rabbbitmq][publishMessage][ERROR]: %s", err.Error())
				failOnError(err, "Failed to publish a message")
			}
			logger.Logger.Infof("[rabbbitmq][publishMessage][INFO] User sucessfully published in queue - MessageID: %s", msg.MessageID)
		}
	}
}

func publishMessage(ctx context.Context, ch *amqp.Channel, p Producer, message *EventDTO) error {
	return ch.PublishWithContext(ctx,
		p.Exchange,
		p.RoutingKey,
		true,
		false,
		amqp.Publishing{
			AppId:       p.Name,
			Timestamp:   time.Now(),
			Body:        message.Body,
			Type:        message.Queue,
			MessageId:   message.MessageID,
			ContentType: "application/json",
		},
	)
}

func failOnError(err error, msg string) {
	if err != nil {
		logger.Logger.Fatalf("%s: %s", msg, err)
	}
}

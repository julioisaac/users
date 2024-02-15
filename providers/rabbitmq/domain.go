package rabbitMQ

type EventDTO struct {
	App         string  `json:"app"`
	Body        []byte  `json:"body"`
	Queue       string  `json:"queue"`
	Priority    byte    `json:"priority"`
	MessageID   string  `json:"messageId"`
	Operation   string  `json:"operation"`
	RetryCount  int     `json:"retryCount"`
	DeliveryTag float64 `json:"deliveryTag"`
}

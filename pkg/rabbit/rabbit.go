package rabbit

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Mailer interface {
	PublishMessage(ctx context.Context, contentType string, body []byte) error
	Close() error
}

type Config struct {
	URL          string `yaml:"url"`
	Exchange     string `yaml:"exchange"`
	ExchangeType string `yaml:"exchange_type"`
	Queue        string `yaml:"queue"`
}

type Connection struct {
	Config
	Connection *amqp.Connection
	Channel    *amqp.Channel
	l          *zap.Logger
}

func New(cfg Config, logger *zap.Logger) (*Connection, error) {
	conn := &Connection{
		Config: cfg,
		l:      logger,
	}
	if err := conn.attemptConnect(); err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	if err := conn.Channel.ExchangeDeclare(
		cfg.Exchange,     // name
		cfg.ExchangeType, // type
		true,             // durable
		false,            // auto-delete
		false,            // internal
		false,            // noWait
		nil,              // arguments
	); err != nil {
		return nil, fmt.Errorf("failed to exchange declare: %s", err)
	}

	if _, err := conn.Channel.QueueDeclare(
		cfg.Queue, // name
		true,      // Durable
		false,     // Auto-delete
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	); err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	if err := conn.Channel.QueueBind(cfg.Queue, "", cfg.Exchange, false, nil); err != nil {
		return nil, fmt.Errorf("failed to bind queue to exchange: %w", err)
	}

	return conn, nil
}

func (c *Connection) attemptConnect() error {
	var err error
	for i := 5; i > 0; i-- {
		if err = c.connect(); err == nil {
			break
		}

		c.l.Info("RabbitMQ is trying to connect again", zap.Int("attempts", i), zap.Error(err))
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return fmt.Errorf("attemptConnect: %w", err)
	}

	return nil
}

func (c *Connection) connect() error {
	var err error

	c.Connection, err = amqp.Dial(c.URL)
	if err != nil {
		return fmt.Errorf("amqp.Dial: %w", err)
	}

	c.Channel, err = c.Connection.Channel()
	if err != nil {
		return fmt.Errorf("Connection.Channel: %w", err)
	}

	return nil
}

func (c *Connection) PublishMessage(ctx context.Context, contentType string, body []byte) error {
	err := c.Channel.PublishWithContext(ctx,
		"",
		c.Config.Queue,
		false,
		false,
		amqp.Publishing{
			ContentType: contentType,
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func (c *Connection) Close() error {
	if c.Connection == nil {
		return nil
	}

	if err := c.Channel.Close(); err != nil {
		return fmt.Errorf("failed to close channel: %w", err)
	}

	if err := c.Connection.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}

	return nil
}

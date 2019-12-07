package rabbitmq

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
	"github.com/thomasobenaus/goms/model"
)

type RabbitMQ struct {
	conn   *amqp.Connection
	logger zerolog.Logger

	host string
	port int

	user     string
	password string

	userQueue       amqp.Queue
	addUserConsumer <-chan amqp.Delivery

	// used to persist the users
	userRepo model.UserRepo
}

// Option represents an option for the api
type Option func(rabbit *RabbitMQ)

// WithLogger adds a configured Logger to the auth
func WithLogger(logger zerolog.Logger) Option {
	return func(rabbit *RabbitMQ) {
		rabbit.logger = logger
	}
}

func Host(host string) Option {
	return func(rabbit *RabbitMQ) {
		rabbit.host = host
	}
}

func Port(port int) Option {
	return func(rabbit *RabbitMQ) {
		rabbit.port = port
	}
}

func UserRepo(userRepo model.UserRepo) Option {
	return func(rabbit *RabbitMQ) {
		rabbit.userRepo = userRepo
	}
}

func New(user, password string, options ...Option) (*RabbitMQ, error) {

	rabbitMq := &RabbitMQ{
		user:     user,
		password: password,
		host:     "localhost",
		port:     5672,
	}

	// apply the options
	for _, opt := range options {
		opt(rabbitMq)
	}

	connStr := fmt.Sprintf("amqp://%s:%s@%s:%d/", rabbitMq.user, rabbitMq.password, rabbitMq.host, rabbitMq.port)
	conn, err := amqp.Dial(connStr)

	if err != nil {
		return nil, err
	}

	rabbitMq.logger.Info().Msgf("Connected to %s:%d", rabbitMq.host, rabbitMq.port)
	rabbitMq.conn = conn

	queue, err := createQueue("AddUserQueue", rabbitMq.conn)
	if err != nil {
		return nil, err
	}
	rabbitMq.userQueue = queue
	rabbitMq.logger.Info().Msgf("Created queue '%s'", rabbitMq.userQueue.Name)

	consumer, err := startAddUserConsumer(rabbitMq.userQueue.Name, rabbitMq.conn, rabbitMq.handleAddUserMsg)
	if err != nil {
		return nil, err
	}
	rabbitMq.addUserConsumer = consumer
	rabbitMq.logger.Info().Msgf("AddUser consumer started")

	return rabbitMq, nil
}

func (rmq *RabbitMQ) Close() error {
	rmq.logger.Info().Msg("Shutting down rabbitmq")
	return rmq.conn.Close()
}

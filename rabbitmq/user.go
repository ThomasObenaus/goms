package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
	"github.com/thomasobenaus/goms/model"
)

func createQueue(queueName string, conn *amqp.Connection) (amqp.Queue, error) {

	ch, err := conn.Channel()
	if err != nil {
		return amqp.Queue{}, err
	}

	return ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
}

func (rmq *RabbitMQ) Add(user model.User) error {
	return fmt.Errorf("Not Implemented yet")
}

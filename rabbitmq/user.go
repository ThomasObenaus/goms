package rabbitmq

import (
	"encoding/json"

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

	ch, err := rmq.conn.Channel()
	if err != nil {
		return err
	}

	body, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",                 // exchange
		rmq.userQueue.Name, // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})

	return err
}

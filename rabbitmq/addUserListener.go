package rabbitmq

import (
	"encoding/json"

	"github.com/streadway/amqp"
	"github.com/thomasobenaus/goms/model"
)

type addUserMsgHandler func(msgBody []byte)

func startAddUserConsumer(queueName string, conn *amqp.Connection, handler addUserMsgHandler) (<-chan amqp.Delivery, error) {

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)

	go func() {
		for data := range msgs {
			handler(data.Body)
		}
	}()

	return msgs, nil
}

func (rmq *RabbitMQ) handleAddUserMsg(msgBody []byte) {
	rmq.logger.Info().Msgf("Received a message: %s", msgBody)

	if rmq.userRepo == nil {
		rmq.logger.Warn().Msgf("User not persisted, UserRepo is nil")
		return
	}

	user := model.User{}
	if err := json.Unmarshal(msgBody, &user); err != nil {
		rmq.logger.Error().Err(err).Msgf("Failed adding user")
	}

	if err := rmq.userRepo.Add(user); err != nil {
		rmq.logger.Error().Err(err).Msgf("Failed adding user %v", user)
		return
	}

	rmq.logger.Info().Msgf("User %v persisted in DB", user)
}

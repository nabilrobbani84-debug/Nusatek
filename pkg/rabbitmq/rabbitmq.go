package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

func ConnectRabbitMQ(url string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	// Declare a queue for property events
	_, err = ch.QueueDeclare(
		"property_events", // name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		return nil, nil, err
	}

	log.Println("Successfully connected to RabbitMQ")
	return conn, ch, nil
}

func PublishEvent(ch *amqp.Channel, queueName string, body []byte) error {
	return ch.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}

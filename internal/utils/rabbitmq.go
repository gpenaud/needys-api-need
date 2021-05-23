package utils

import (
  "fmt"
  "log"
	"github.com/streadway/amqp"
  // local packages
  "github.com/gpenaud/needys-api-need/internal/models"
)

func SendAmqpMessages(message []byte) {
  rabbitmq_connection_parameters := fmt.Sprintf(
    "amqp://%s:%s@%s:%s/",
    models.Cfg.Rabbitmq.Username,
    models.Cfg.Rabbitmq.Password,
    models.Cfg.Rabbitmq.Host,
    models.Cfg.Rabbitmq.Port,
  )

  conn, err := amqp.Dial(rabbitmq_connection_parameters)

	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	FailOnError(err, "Failed to declare a queue")

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	FailOnError(err, "Failed to publish a message")

  log.Printf(" [x] Sent %s", message)
}

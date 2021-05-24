package rabbitmq_sender

import (
  "fmt"
	"github.com/streadway/amqp"
  // local packages
  // "github.com/gpenaud/needys-api-need/internal/models"
  "github.com/gpenaud/needys-api-need/internal/config"
  "github.com/gpenaud/needys-api-need/pkg/log"
)

func SendAmqpMessages(message []byte) {
  rabbitmq_connection_parameters := fmt.Sprintf(
    "amqp://%s:%s@%s:%s/",
    config.Cfg.Rabbitmq.Username,
    config.Cfg.Rabbitmq.Password,
    config.Cfg.Rabbitmq.Host,
    config.Cfg.Rabbitmq.Port,
  )

  conn, err := amqp.Dial(rabbitmq_connection_parameters)
  if err != nil {
    log.ErrorLogger.Fatalln("Failed to close RabbitMQ")
  }

	defer conn.Close()

	ch, err := conn.Channel()
  if err != nil {
    log.ErrorLogger.Fatalln("Failed to close the channel")
  }

	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

  if err != nil {
    log.ErrorLogger.Fatalln("Failed to declare a queue")
  }

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
  )

  if err != nil {
    log.ErrorLogger.Fatalln("Failed to publish a message")
  }

  log.InfoLogger.Printf(" [x] Sent %s", message)
}

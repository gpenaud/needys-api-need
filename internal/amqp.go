package internal

import (
  amqp "github.com/streadway/amqp"
  fmt  "fmt"
  log  "github.com/sirupsen/logrus"
)

var amqpLog *log.Entry

func init() {
  amqpLog = log.WithFields(log.Fields{
    "_file": "internal/amqp.go",
    "_type": "messaging",
  })
}

func (a* Application) SendAmqpMessages(message []byte) {
  rabbitmq_connection_parameters := fmt.Sprintf(
    "amqp://%s:%s@%s:%s/",
    a.Config.Rabbitmq.Username,
    a.Config.Rabbitmq.Password,
    a.Config.Rabbitmq.Host,
    a.Config.Rabbitmq.Port,
  )

  conn, err := amqp.Dial(rabbitmq_connection_parameters)
  if err != nil {
    amqpLog.Fatal("Failed to close RabbitMQ")
  }

	defer conn.Close()

	ch, err := conn.Channel()
  if err != nil {
    amqpLog.Fatal("Failed to close the channel")
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
    amqpLog.Fatal("Failed to declare a queue")
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
    amqpLog.Fatal("Failed to publish a message")
  }

  amqpLog.Info(" [x] Sent %s", message)
}

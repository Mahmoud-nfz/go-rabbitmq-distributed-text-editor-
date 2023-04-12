package rabbitmq

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Lock(queue_name string, content string) {
	conn := establish_connection()
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")

	err = ch.ExchangeDeclare("Lock 1", "fanout", true, false, false, false, nil)
	FailOnError(err, "Failed to declare exchange Lock 1")
	err = ch.ExchangeDeclare("Lock 2", "fanout", true, false, false, false, nil)
	FailOnError(err, "Failed to declare exchange Lock 2")


	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()


	err = ch.PublishWithContext(ctx,
		queue_name, // exchange
		"",         // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(content),
		})

	FailOnError(err, "Failed to Lock a text field")

	log.Printf(" [x] Locked %s for task: %s", content, queue_name)

}

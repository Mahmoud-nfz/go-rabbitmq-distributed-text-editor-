package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func establish_connection() *amqp.Connection {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")

	return conn
}

func declare_exchanges(ch *amqp.Channel) {
	err := ch.ExchangeDeclare("display", "fanout", true, false, false, false, nil)
	FailOnError(err, "Failed to declare display exchange")
	err = ch.ExchangeDeclare("tasks", "direct", true, false, false, false, nil)
	FailOnError(err, "Failed to declare tasks exchange")
}

func declare_queues(ch *amqp.Channel, id int) amqp.Queue {
	my_queue := fmt.Sprintf("display_%d", id)
	q, err := ch.QueueDeclare(my_queue, true, false, false, false, nil)
	FailOnError(err, "Failed to declare queue")
	err = ch.QueueBind(q.Name, "", "display", false, nil)
	FailOnError(err, "Failed to bind queue")

	return q
}

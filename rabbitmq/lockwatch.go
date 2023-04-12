package rabbitmq

import (
	// "fmt"
	"log"

	// "math/rand"

	"fyne.io/fyne/v2/data/binding"
	amqp "github.com/rabbitmq/amqp091-go"
)

func LockWatch(task1 *binding.String, task2 *binding.String, rand_int string) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")

	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")

	err = ch.ExchangeDeclare("Lock 1", "fanout", true, false, false, false, nil)
	FailOnError(err, "Failed to declare display exchange task 1")

	err = ch.ExchangeDeclare("Lock 2", "fanout", true, false, false, false, nil)
	FailOnError(err, "Failed to declare display exchange task 2")

	q1, err := ch.QueueDeclare("Lock 1"+rand_int, false, false, false, false, nil)
	FailOnError(err, "Failed to declare queue task 1")
	q2, err := ch.QueueDeclare("Lock 2"+rand_int, false, false, false, false, nil)
	FailOnError(err, "Failed to declare queue task 2")

	// fmt.Println("My queue name is: ", rand_int)

	err = ch.QueueBind(q1.Name, "", "Lock 1", false, nil)
	FailOnError(err, "Failed to bind queue task 1")
	err = ch.QueueBind(q2.Name, "", "Lock 2", false, nil)
	FailOnError(err, "Failed to bind queue task 2")

	msgs_task1, err := ch.Consume(q1.Name, "", true, false, false, false, nil)
	FailOnError(err, "Failed to consume queue task 1")
	msgs_task2, err := ch.Consume(q2.Name, "", true, false, false, false, nil)
	FailOnError(err, "Failed to consume queue task 2")

	go func() {
		for d := range msgs_task1 {
			log.Printf("Received a lock notification for task 1: %s", d.Body)
			// current, err := (*task1).Get()
			// FailOnError(err, "Failed to consume queue task 1")
			(*task1).Set("(Locked by " + string(d.Body) + ")")
		}
	}()

	go func() {
		for d := range msgs_task2 {
			log.Printf("Received a lock notification for task 1: %s", d.Body)
			// current, err := (*task2).Get()
			// FailOnError(err, "Failed to consume queue task 2")
			(*task2).Set("(Locked by " + string(d.Body) + ")")
			
		}

	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	<-make(chan bool)

}

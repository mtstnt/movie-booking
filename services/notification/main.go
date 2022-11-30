package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	DEBUG_MAILHOST = "mailhog:1025"
)

func InitializeQueueConnection(username, password, host string, port int) (
	*amqp.Connection, *amqp.Channel, *amqp.Queue, error,
) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d", username, password, host, port))
	if err != nil {
		return nil, nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, nil, err
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare("notification_queue", true, false, false, false, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	return conn, ch, &queue, nil
}

func run() error {
	conn, ch, q, err := InitializeQueueConnection(
		"guest", "guest", "rabbitmq", 5672,
	)
	if err != nil {
		return err
	}

	defer conn.Close()
	defer ch.Close()

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	var forever chan struct{}

	go func() {
		for msg := range msgs {
			from := "test@email.com"
			to := []string{"receiver1@email.com", "receiver2@email.com"}
			subject := "<h1>Hello world!</h1>"

			if err := msg.Ack(false); err != nil {
				log.Printf("Error acknowledging message %s.\n", msg.MessageId)
			}

			var message struct {
				To      string
				Subject string
			}
			if err := json.Unmarshal(msg.Body, &message); err != nil {
				log.Printf("Error reading request in correct format: \n%s\n", string(msg.Body))
			}

			if err := SendNotifications(from, to, subject); err != nil {
				log.Printf("Error sending email to: %v, subject: %s\n", to, subject)
			}
		}
	}()

	log.Println("Waiting for messages...")
	<-forever
	return nil
}

func SendNotifications(from string, to []string, subject string) error {
	if err := smtp.SendMail(DEBUG_MAILHOST, nil, from, to, []byte(subject)); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

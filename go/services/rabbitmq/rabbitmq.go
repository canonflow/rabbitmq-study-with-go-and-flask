package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

var RabbitMQClient *RabbitMQ

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func Test() {
	fmt.Println(os.Getenv("RABBITMQ_CONNECTION_URL"))
}

func NewRabbitMQConnection() {
	// TODO: Connect to RabbitMQ
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_CONNECTION_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	// TODO: Open a RabbitMQ Channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a RabbitMQ channel: %s", err)
	}

	RabbitMQClient = &RabbitMQ{
		Conn:    conn,
		Channel: ch,
	}
}

func (r *RabbitMQ) PublishMessage(message map[string]string, queue string) error {
	// TODO: Declare the queue to ensure it exists
	q, err := r.Channel.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
		return err
	}

	// TODO: Convert the message to JSON Format
	body, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Failed to marshal message: %v", err)
		return err
	}

	// TODO: Publish the message to the queue
	err = r.Channel.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
		return err
	}

	log.Printf("Message has been send to RabbitMQ queue: %s", message)

	return nil
}

func (r *RabbitMQ) ConsumeQueue(queue string) (<-chan amqp.Delivery, error) {
	// TODO: Declare the queue to ensure it exists
	q, err := r.Channel.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
		return nil, err
	}

	// SUBSCRIBE TO THE QUEUE
	msgs, err := r.Channel.Consume(
		q.Name, // queue
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Failed to register a RabbitMQ consumer: %s", err)
		return nil, err
	}

	return msgs, nil
}

func (r *RabbitMQ) CloseConnection() {
	r.Channel.Close()
	r.Conn.Close()
}

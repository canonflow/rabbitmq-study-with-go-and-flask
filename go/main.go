package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/canonflow/rabbitmq-study-with-go-and-flask/controllers"
	"github.com/canonflow/rabbitmq-study-with-go-and-flask/services/rabbitmq"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %+v", err)
	}

	// TODO: Init RabbitMQClient
	rabbitmq.NewRabbitMQConnection()
}

func main() {
	defer rabbitmq.RabbitMQClient.CloseConnection()

	// TODO: Create Endpoint
	gin := gin.Default()

	gin.POST("/publish", controllers.SendMessage)

	// TODO: Consume RabbitMQ
	msgs, err := rabbitmq.RabbitMQClient.ConsumeQueue("second_queue")
	if err != nil {
		log.Fatalf("Failed to consume RabbitMQ queue: %s", err)
		return
	}

	go func() {
		for d := range msgs {
			var publishedMessage map[string]string

			err := json.Unmarshal(d.Body, &publishedMessage)
			if err != nil {
				log.Printf("Error reading coffee order (please check the JSON format): %s", err)
				continue
			}

			fmt.Println("Received Message: ", publishedMessage)
		}
	}()

	gin.Run(":" + "8990")
}

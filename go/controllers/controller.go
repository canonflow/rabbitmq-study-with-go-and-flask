package controllers

import (
	"net/http"

	"github.com/canonflow/rabbitmq-study-with-go-and-flask/services/rabbitmq"
	"github.com/gin-gonic/gin"
)

func SendMessage(c *gin.Context) {
	// TODO: Send the message to RabbitMQ
	queueName := "first_queue"
	message := map[string]string{
		"message": "Hello From Golang",
		"state":   "Publish to RabbitMQ",
	}

	err := rabbitmq.RabbitMQClient.PublishMessage(message, queueName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to send coffee order to RabbitMQ queue"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Coffee order has been sent to RabbitMQ queue."})
}

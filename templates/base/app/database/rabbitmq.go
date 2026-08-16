package database

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
)

func GetRabbitMQConnection() *amqp.Connection {

	host := os.Getenv("RABBITMQ_HOST")
	user := os.Getenv("RABBITMQ_USER")
	pass := os.Getenv("RABBITMQ_PASS")
	port := os.Getenv("RABBITMQ_PORT")
	vhost := os.Getenv("RABBITMQ_VHOST")

	uri := fmt.Sprintf("amqp://%s:%s@%s:%s/%s", user, pass, host, port, vhost)

	conn, err := amqp.Dial(uri)

	if err != nil {

		log.Printf("got error connecting to rabbitMQ %s with %s", err.Error(), uri)
		return nil
	}

	return conn
}

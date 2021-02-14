package rabbit

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

// Client - rabbitmq client
type Client struct {
	connection *amqp.Connection
}

// CreateClient - generate a rabbitmq client
func CreateClient(connectionURL string) (*Client, error) {
	conn, err := amqp.Dial(connectionURL)
	client := Client{connection: conn}
	log.Printf("[rabbit] client connected successfully")
	return &client, err
}

// Publish - publish data to queue
func (client *Client) Publish(routingKey string, payload interface{}) error {
	ch, err := client.connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		routingKey, // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return err
	}
	data, _ := json.Marshal(payload)
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		})
	log.Printf("[rabbit] published %s on %s", data, q.Name)
	return err
}

// Close - Close the client
func (client *Client) Close() {
	client.connection.Close()
}

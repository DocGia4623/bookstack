package messaging

import (
	"bookstack/config"
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQ(conf *config.Config) (*RabbitMQ, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/",
		conf.RabbitMqUser,
		conf.RabbitMQPassword,
		conf.RabbitMQHost,
		conf.RabbitMQPort))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %v", err)
	}

	return &RabbitMQ{
		conn:    conn,
		channel: ch,
	}, nil
}

func (r *RabbitMQ) PublishNewOrder(orderID uint, address string) error {
	queue, err := r.channel.QueueDeclare(
		"new_orders", // queue name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return err
	}

	message := map[string]interface{}{
		"order_id": orderID,
		"address":  address,
	}

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = r.channel.PublishWithContext(
		context.Background(),
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	return err
}

func (r *RabbitMQ) ConsumeNewOrders(handler func(orderID uint, address string)) error {
	queue, err := r.channel.QueueDeclare(
		"new_orders", // queue name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return err
	}

	msgs, err := r.channel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			var data map[string]interface{}
			if err := json.Unmarshal(msg.Body, &data); err != nil {
				log.Printf("Error decoding message: %v", err)
				continue
			}

			orderID := uint(data["order_id"].(float64))
			address := data["address"].(string)
			handler(orderID, address)
		}
	}()

	return nil
}

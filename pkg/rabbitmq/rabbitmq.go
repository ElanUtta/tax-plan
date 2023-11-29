package rabbitmq

import (
	"context"
	"encoding/json"
	"time"

	"github.com/elanutta/go-intensivo/internal/usecase"
	amqp "github.com/rabbitmq/amqp091-go"
)

func OpenChannel() (*amqp.Channel, error) {
	// Open connection with rabbitmq
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func Consume(ch *amqp.Channel, out chan amqp.Delivery) error {
	msgs, err := ch.Consume(
		"order",
		"go-consumer",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for msg := range msgs {
		out <- msg
	}
	return nil
}

func Publisher(ch *amqp.Channel, text_msg usecase.OrderInput) error {

	err := ch.ExchangeDeclare(
		"oerder_exchange", // name
		"fanout",          // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	)

	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	json_message, err := json.Marshal(text_msg)
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(ctx,
		"oerder_exchange", // exchange
		"nao_sei",         // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(json_message),
		})

	ch.QueueBind(
		"order",
		"nao_sei",
		"oerder_exchange",
		false,
		nil,
	)

	return nil

	// msg := amqp.Publishing{
	// 	DeliveryMode: amqp.Persistent,
	// 	Timestamp:    time.Now(),
	// 	ContentType:  "application/json",
	// 	Body:         []byte(json_message),
	// }

	// err = ch.PublishWithContext(
	// 	ctx,
	// 	"order",
	// 	"go-delivery",
	// 	false,
	// 	false,
	// 	msg,
	// )
	// println("foi publicado a messagem ", json_message)

	// if err != nil {
	// 	// Since publish is asynchronous this can happen if the network connection
	// 	// is reset or if the server has run out of resources.
	// 	log.Fatalf("basic.publish: %v", err)
	// }

	// return nil

}

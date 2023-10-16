package rmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	conn         *amqp.Connection
	ch           *amqp.Channel
	exchangeName string
	routingKey   string
}

func NewProducer(ops MQOptions) (*Producer, error) {
	conn, ch, err := connect(ops)
	if err != nil {
		return nil, err
	}

	return &Producer{
		conn:         conn,
		ch:           ch,
		exchangeName: ops.Ename,
		routingKey:   ops.RoutingKey,
	}, nil
}

func (p *Producer) Publish(ctx context.Context, body []byte) error {
	return p.ch.PublishWithContext(ctx,
		p.exchangeName,
		p.routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (p *Producer) Close() error {
	return p.conn.Close()
}

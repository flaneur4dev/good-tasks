package rmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	cs "github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/contracts"
)

type Consumer struct {
	conn        *amqp.Connection
	ch          *amqp.Channel
	queueName   string
	consumerTag string
}

func NewConsumer(ops MQOptions, tag string) (*Consumer, error) {
	conn, ch, err := connect(ops)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		conn:        conn,
		ch:          ch,
		queueName:   ops.Qname,
		consumerTag: tag,
	}, nil
}

func (c *Consumer) Consume(ctx context.Context) (<-chan cs.NotificationMessage, error) {
	deliveries, err := c.ch.ConsumeWithContext(ctx, c.queueName, c.consumerTag, false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	messages := make(chan cs.NotificationMessage)
	go func() {
		defer close(messages)

		for d := range deliveries {
			if err := d.Ack(false); err != nil {
				fmt.Println(err)
			}

			msg := cs.NotificationMessage{
				ID:              d.MessageId,
				ContentType:     d.ContentType,
				ContentEncoding: d.ContentEncoding,
				Body:            d.Body,
			}

			messages <- msg
		}
	}()

	return messages, nil
}

func (c *Consumer) Close() error {
	err := c.ch.Cancel(c.consumerTag, false)
	if err != nil {
		return err
	}

	return c.conn.Close()
}

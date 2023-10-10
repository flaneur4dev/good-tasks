package rmq

import amqp "github.com/rabbitmq/amqp091-go"

type MQOptions struct {
	URL         string
	RoutingKey  string
	Ename       string
	Etype       string
	Edurable    bool
	EautoDelete bool
	Einternal   bool
	EnoWait     bool
	Qname       string
	Qdurable    bool
	QautoDelete bool
	Qexclusive  bool
	QnoWait     bool
}

func connect(ops MQOptions) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(ops.URL)
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	err = ch.ExchangeDeclare(ops.Ename, ops.Etype, ops.Edurable, ops.EautoDelete, ops.Einternal, ops.EnoWait, nil)
	if err != nil {
		return nil, nil, err
	}

	q, err := ch.QueueDeclare(ops.Qname, ops.Qdurable, ops.QautoDelete, ops.Qexclusive, ops.QnoWait, nil)
	if err != nil {
		return nil, nil, err
	}

	err = ch.QueueBind(q.Name, ops.RoutingKey, ops.Ename, false, nil)
	if err != nil {
		return nil, nil, err
	}

	return conn, ch, nil
}

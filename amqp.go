package queue

import "github.com/streadway/amqp"

// AMQPBackend provides AMQP-based backend to manage queues.
// https://en.wikipedia.org/wiki/Advanced_Message_Queuing_Protocol
// Suitable for multi-host, multi-process and multithreaded environment
type AMQPBackend struct {
	conn  *amqp.Connection
	codec Codec
}

// NewAMQPBackend creates new AMQPBackend
func NewAMQPBackend(url string) (*AMQPBackend, error) {
	conn, err := amqp.Dial(url)

	if err != nil {
		return nil, err
	}

	b := &AMQPBackend{conn: conn}

	return b.Codec(NewGOBCodec()), nil
}

// Codec sets codec to encode/decode objects in queues. GOBCodec is default.
func (b *AMQPBackend) Codec(c Codec) *AMQPBackend {
	b.codec = c
	return b
}

// Put adds value to the end of a queue.
func (b *AMQPBackend) Put(queueName string, value interface{}) error {
	data, err := b.codec.Marshal(value)

	if err != nil {
		return err
	}

	ch, err := b.conn.Channel()

	if err != nil {
		return err
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		true,      // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		return err
	}

	return ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/octet-stream",
			Body:        data,
		})
}

// Get removes the first element from a queue and put it in the value pointed to by v
func (b *AMQPBackend) Get(queueName string, v interface{}) error {
	ch, err := b.conn.Channel()

	if err != nil {
		return err
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		true,      // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		return err
	}

	d := <-msgs

	return b.codec.Unmarshal(d.Body, v)
}

func (b *AMQPBackend) RemoveQueue(queueName string) error {
	return nil
}

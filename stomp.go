package queue

import (
	"github.com/go-stomp/stomp"
)

// StompBackend provides stomp-based backend to manage queues.
// Suitable for multi-host, multi-process and multithreaded environment
type StompBackend struct {
	conn  *stomp.Conn
	codec Codec
}

// NewStompBackend creates new RedisBackend
func NewStompBackend(addr string, opts ...func(*stomp.Conn) error) (*StompBackend, error) {
	conn, err := stomp.Dial("tcp", addr, opts...)

	if err != nil {
		return nil, err
	}

	b := &StompBackend{conn: conn}
	return b.Codec(NewGOBCodec()), nil
}

// Codec sets codec to encode/decode objects in queues. GOBCodec is default.
func (b *StompBackend) Codec(c Codec) *StompBackend {
	b.codec = c
	return b
}

// Put adds value to the end of a queue.
func (b *StompBackend) Put(queueName string, value interface{}) error {
	data, err := b.codec.Marshal(value)

	if err != nil {
		return err
	}

	return b.conn.Send(queueName, "application/octet-stream", data, stomp.SendOpt.Receipt)
}

// Get removes the first element from a queue and put it in the value pointed to by v
func (b *StompBackend) Get(queueName string, value interface{}) error {
	sub, err := b.conn.Subscribe(queueName, stomp.AckClientIndividual)

	if err != nil {
		return err
	}

	defer sub.Unsubscribe()
	msg := <-sub.C
	err = b.conn.Ack(msg)

	if err != nil {
		return err
	}

	return b.codec.Unmarshal(msg.Body, value)
}

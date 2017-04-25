package queue

import (
	"testing"
	"github.com/stretchr/testify/suite"
	"fmt"
)

type AMQPJSONTestSuite struct {
	baseSuite
}

func (suite *AMQPJSONTestSuite) SetupTest() {
	b, err := NewAMQPBackend("amqp://guest:guest@192.168.99.100:5672")

	if err != nil {
		fmt.Printf("amqp err: %v\n", err)
	}

	b.Codec(NewJSONCodec())
	Init(b)
}

func TestAMQPJSONTestSuite(t *testing.T) {
	suite.Run(t, new(AMQPJSONTestSuite))
}

type AMQPGOBTestSuite struct {
	baseSuite
}

func (suite *AMQPGOBTestSuite) SetupTest() {
	b, err := NewAMQPBackend("amqp://guest:guest@192.168.99.100:5672")

	if err != nil {
		fmt.Printf("amqp err: %v\n", err)
	}

	Init(b)
}

func TestAMQPGOBTestSuite(t *testing.T) {
	suite.Run(t, new(AMQPGOBTestSuite))
}

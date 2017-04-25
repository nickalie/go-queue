package queue

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type AMQPJSONTestSuite struct {
	baseSuite
}

func (suite *AMQPJSONTestSuite) SetupTest() {
	host, _ := os.LookupEnv("TESTS_HOST")
	b, err := NewAMQPBackend("amqp://guest:guest@" + host + ":5672")

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
	host, _ := os.LookupEnv("TESTS_HOST")
	b, err := NewAMQPBackend("amqp://guest:guest@" + host + ":5672")

	if err != nil {
		fmt.Printf("amqp err: %v\n", err)
	}

	Init(b)
}

func TestAMQPGOBTestSuite(t *testing.T) {
	suite.Run(t, new(AMQPGOBTestSuite))
}

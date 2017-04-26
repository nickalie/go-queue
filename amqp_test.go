package queue

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
)

type AMQPTestSuite struct {
	baseSuite
}

func (suite *AMQPTestSuite) SetupTest() {
	host, _ := os.LookupEnv("TESTS_HOST")
	b, err := NewAMQPBackend("amqp://guest:guest@" + host + ":5672")

	if err != nil {
		fmt.Printf("amqp err: %v\n", err)
	}

	Use(b)
}

func TestAMQPTestSuite(t *testing.T) {
	suite.Run(t, new(AMQPTestSuite))
}

func TestInvalidAMQPUrl(t *testing.T) {
	b, err := NewAMQPBackend("https://google.com")
	assert.NotNil(t, err)
	assert.Nil(t, b)
}

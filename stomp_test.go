package queue

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type StompTestSuite struct {
	baseSuite
}

func (suite *StompTestSuite) SetupTest() {
	host, _ := os.LookupEnv("TESTS_HOST")
	b, err := NewStompBackend(host + ":61613")

	if err != nil {
		fmt.Printf("stomp err: %v\n", err)
	}

	Use(b)
}

func TestStompTestSuite(t *testing.T) {
	//TODO find a way to run stomp on circleci
	if _, ok := os.LookupEnv("CIRCLECI"); !ok {
		suite.Run(t, new(StompTestSuite))
	}
}

func TestInvalidStompUrl(t *testing.T) {
	b, err := NewStompBackend("https://google.com")
	assert.NotNil(t, err)
	assert.Nil(t, b)
}

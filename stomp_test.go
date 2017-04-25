package queue

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type StompJSONTestSuite struct {
	baseSuite
}

func (suite *StompJSONTestSuite) SetupTest() {
	host, _ := os.LookupEnv("TESTS_HOST")
	b, err := NewStompBackend(host + ":61613")

	if err != nil {
		fmt.Printf("stomp err: %v\n", err)
	}

	b.Codec(NewJSONCodec())
	Init(b)
}

func TestStompJSONTestSuite(t *testing.T) {
	suite.Run(t, new(StompJSONTestSuite))
}

type StompGOBTestSuite struct {
	baseSuite
}

func (suite *StompGOBTestSuite) SetupTest() {
	host, _ := os.LookupEnv("TESTS_HOST")
	b, err := NewStompBackend(host + ":61613")

	if err != nil {
		fmt.Printf("stomp err: %v\n", err)
	}

	Init(b)
}

func TestStompGOBTestSuite(t *testing.T) {
	suite.Run(t, new(StompGOBTestSuite))
}

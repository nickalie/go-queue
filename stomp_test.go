package queue

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

type StompJSONTestSuite struct {
	baseSuite
}

func (suite *StompJSONTestSuite) SetupTest() {
	b, err := NewStompBackend("192.168.99.100:61613")

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
	b, err := NewStompBackend("192.168.99.100:61613")

	if err != nil {
		fmt.Printf("stomp err: %v\n", err)
	}

	Init(b)
}

func TestStompGOBTestSuite(t *testing.T) {
	suite.Run(t, new(StompGOBTestSuite))
}

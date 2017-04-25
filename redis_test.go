package queue

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

type RedisJSONTestSuite struct {
	baseSuite
}

func (suite *RedisJSONTestSuite) SetupTest() {
	b, err := NewRedisBackend("redis://localhost:6379")

	if err != nil {
		fmt.Printf("redis err: %v\n", err)
	}

	b.Codec(NewJSONCodec())
	Init(b)
}

func TestRedisJSONTestSuite(t *testing.T) {
	suite.Run(t, new(RedisJSONTestSuite))
}

type RedisGOBTestSuite struct {
	baseSuite
}

func (suite *RedisGOBTestSuite) SetupTest() {
	b, err := NewRedisBackend("redis://localhost:6379")

	if err != nil {
		fmt.Printf("redis err: %v\n", err)
	}

	Init(b)
}

func TestRedisGOBTestSuite(t *testing.T) {
	suite.Run(t, new(RedisGOBTestSuite))
}

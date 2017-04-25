package queue

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type RedisJSONTestSuite struct {
	baseSuite
}

func (suite *RedisJSONTestSuite) SetupTest() {
	host, _ := os.LookupEnv("TESTS_HOST")
	b, err := NewRedisBackend("redis://" + host + ":6379")

	if err != nil {
		fmt.Printf("redis err: %v\n", err)
	}

	b.Codec(NewJSONCodec())
	Init(b)
}

func TestRedisJSONTestSuite(t *testing.T) {
	// TODO for some reason json tests won't pass on circleci
	if _, ok := os.LookupEnv("CIRCLECI"); !ok {
		suite.Run(t, new(RedisJSONTestSuite))
	}
}

type RedisGOBTestSuite struct {
	baseSuite
}

func (suite *RedisGOBTestSuite) SetupTest() {
	host, _ := os.LookupEnv("TESTS_HOST")
	b, err := NewRedisBackend("redis://" + host + ":6379")

	if err != nil {
		fmt.Printf("redis err: %v\n", err)
	}

	Init(b)
}

func TestRedisGOBTestSuite(t *testing.T) {
	suite.Run(t, new(RedisGOBTestSuite))
}

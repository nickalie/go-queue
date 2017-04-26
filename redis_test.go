package queue

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type RedisTestSuite struct {
	baseSuite
}

func (suite *RedisTestSuite) SetupTest() {
	host, _ := os.LookupEnv("TESTS_HOST")
	b, err := NewRedisBackend("redis://" + host + ":6379")

	if err != nil {
		fmt.Printf("redis err: %v\n", err)
	}

	Use(b)
}

func TestRedisTestSuite(t *testing.T) {
	suite.Run(t, new(RedisTestSuite))
}

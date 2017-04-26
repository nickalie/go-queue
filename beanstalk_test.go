package queue

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type BeanstalkTestSuite struct {
	baseSuite
}

func (suite *BeanstalkTestSuite) SetupTest() {
	host, _ := os.LookupEnv("TESTS_HOST")
	b, err := NewBeanstalkBackend(host + ":11300")

	if err != nil {
		fmt.Printf("beanstalk err: %v\n", err)
	}

	Use(b)
}

func TestBeanstalkTestSuite(t *testing.T) {
	suite.Run(t, new(BeanstalkTestSuite))
}

func TestInvalidBeanstalkUrl(t *testing.T) {
	b, _ := NewBeanstalkBackend("google.com")
	assert.NotNil(t, b)
	err := b.Put(randString(10), randUser())
	assert.NotNil(t, err)
	err = b.Get(randString(10), &testUser{})
	assert.NotNil(t, err)
}

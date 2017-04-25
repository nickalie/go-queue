package queue

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type FSJSONTestSuite struct {
	baseSuite
}

func (suite *FSJSONTestSuite) SetupTest() {
	os.RemoveAll("data")
	b, err := NewFSBackend("data")

	if err != nil {
		fmt.Printf("fs err: %v\n", err)
	}

	b.Codec(NewJSONCodec())
	Init(b)
}

func TestFSJSONTestSuite(t *testing.T) {
	suite.Run(t, new(FSJSONTestSuite))
}

type FSGOBTestSuite struct {
	baseSuite
}

func (suite *FSGOBTestSuite) SetupTest() {
	os.RemoveAll("data")
	b, err := NewFSBackend("data")

	if err != nil {
		fmt.Printf("fs err: %v\n", err)
	}

	Init(b)
}

func TestFSGOBTestSuite(t *testing.T) {
	suite.Run(t, new(FSGOBTestSuite))
}

package queue

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type FSTestSuite struct {
	baseSuite
}

func (suite *FSTestSuite) SetupTest() {
	os.RemoveAll("data")
	b, err := NewFSBackend("data")

	if err != nil {
		fmt.Printf("fs err: %v\n", err)
	}

	Use(b)
}

func TestFSTestSuite(t *testing.T) {
	suite.Run(t, new(FSTestSuite))
}

package queue

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type ChannelTestSuite struct {
	baseSuite
}

func (suite *ChannelTestSuite) SetupTest() {
	Init(NewChannelBackend())
}

func TestChannelTestSuite(t *testing.T) {
	suite.Run(t, new(ChannelTestSuite))
}

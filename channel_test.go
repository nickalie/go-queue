package queue

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type ChannelTestSuite struct {
	baseSuite
}

func (suite *ChannelTestSuite) SetupTest() {
	Use(NewChannelBackend())
}

func TestChannelTestSuite(t *testing.T) {
	suite.Run(t, new(ChannelTestSuite))
}

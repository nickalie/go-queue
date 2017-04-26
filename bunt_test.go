package queue

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"log"
	"os"
	"testing"
)

type BuntTestSuite struct {
	baseSuite
	b *BuntBackend
}

func (suite *BuntTestSuite) SetupTest() {

	if suite.b != nil {
		err := suite.b.Close()

		if err != nil {
			log.Fatal(err)
		}
	}

	os.Remove("data.db")

	b, err := NewBuntBackend("data.db")

	if err != nil {
		fmt.Printf("bunt err: %v\n", err)
	}

	suite.b = b

	Use(b)
}

func TestBuntTestSuite(t *testing.T) {
	suite.Run(t, new(BuntTestSuite))
}

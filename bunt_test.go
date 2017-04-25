package queue

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"log"
	"os"
	"testing"
)

type BuntJSONTestSuite struct {
	baseSuite
	b *BuntBackend
}

func (suite *BuntJSONTestSuite) SetupTest() {

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

	b.Codec(NewJSONCodec())
	Init(b)
	suite.b = b
}

func TestBuntJSONTestSuite(t *testing.T) {
	suite.Run(t, new(BuntJSONTestSuite))
}

type BuntGOBTestSuite struct {
	baseSuite
	b *BuntBackend
}

func (suite *BuntGOBTestSuite) SetupTest() {

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

	Init(b)
}

func TestBuntGOBTestSuite(t *testing.T) {
	suite.Run(t, new(BuntGOBTestSuite))
}

package queue

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewJSONCodec(t *testing.T) {
	c := NewJSONCodec()
	expected := randMap(100)
	js, err := c.Marshal(&expected)
	assert.Nil(t, err)
	actual := make(map[string]interface{})
	err = c.Unmarshal(js, &actual)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

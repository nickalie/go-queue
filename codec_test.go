package queue

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewJSONCodec(t *testing.T) {
	c := NewJSONCodec()
	expected := randUser()
	js, err := c.Marshal(&expected)
	assert.Nil(t, err)
	var actual testUser
	err = c.Unmarshal(js, &actual)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

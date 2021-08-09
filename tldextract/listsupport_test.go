package tldextract

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CreateList(t *testing.T) {
	assert := assert.New(t)

	actual := CreateList()

	assert.NotNil(actual, "Not Nil")
	assert.NotEqual(0, len(actual), "Length is not 0")
}

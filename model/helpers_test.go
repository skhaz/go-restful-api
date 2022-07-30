package model

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLast(t *testing.T) {
	expected := 3
	actual, ok := Last([]int{1, 2, 3})
	assert.True(t, ok)
	assert.Equal(t, expected, actual)
}

func TestGetLastEmpty(t *testing.T) {
	expected := 0
	actual, ok := Last([]int{})
	assert.False(t, ok)
	assert.Equal(t, expected, actual)
}

func TestNoError(t *testing.T) {
	assert.True(t, NoError(nil))
	assert.False(t, NoError(errors.New("an error has occurred")))
}

func TestJSONRemarshalError(t *testing.T) {
	_, err := JSONRemarshal([]byte("{"))
	assert.Error(t, err)
}

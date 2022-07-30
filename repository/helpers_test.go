package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAnyTimeMatch(t *testing.T) {
	assert.True(t, AnyTime{}.Match(time.Now()))
}

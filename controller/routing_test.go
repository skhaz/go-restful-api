package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterRoutes(t *testing.T) {
	s := InitServer()

	s.registerRoutes()

	assert.NotEmpty(t, s.router.Routes())
}

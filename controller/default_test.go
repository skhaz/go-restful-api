package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNoRoute(t *testing.T) {
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)

	NoRoute(ctx)

	assert.Equal(t, r.Code, http.StatusNotFound)
	assert.Equal(t, r.Header().Get("Content-Type"), "application/problem+json")
	assert.Contains(t, r.Body.String(), "errors:http/not-found")
}

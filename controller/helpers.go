package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pmoule/go2hal/hal"
)

var validate = validator.New()

var encoder = hal.NewEncoder()

func WriteHAL(ctx *gin.Context, statusCode int, resource hal.Resource) {
	encoder.WriteTo(ctx.Writer, statusCode, resource)
}

func WriteNoContent(ctx *gin.Context) {
	ctx.JSON(http.StatusNoContent, nil)
}

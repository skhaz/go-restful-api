package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"schneider.vip/problem"
)

func NoRoute(ctx *gin.Context) {
	problem.New(
		problem.Title("Not Found"),
		problem.Type("errors:http/not-found"),
		problem.Status(http.StatusNotFound),
	).WriteTo(ctx.Writer)
}

package controller

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"skhaz.dev/rest/repository"
)

type Server struct {
	router *gin.Engine
}

func (s *Server) registerRoutes() {
	var router = s.router

	router.NoRoute(NoRoute)

	workspaces := router.Group("/workspaces")
	{
		workspaces.GET("", GetWorkspaces)
		workspaces.POST("", CreateWorkspace)
		workspaces.GET("/:uuid", GetWorkspace)
		workspaces.PATCH("/:uuid", UpdateWorkspace)
		workspaces.DELETE("/:uuid", DeleteWorkspace)
	}
}

func InitServer() *Server {
	return &Server{router: gin.Default()}
}

func (s *Server) Run() (err error) {
	s.registerRoutes()

	s.router.Use(gzip.Gzip(gzip.DefaultCompression))

	err = s.router.Run()
	return
}

func (s *Server) SetRepositoryRegistry(rr *repository.RepositoryRegistry) {
	s.router.Use(func(c *gin.Context) {
		c.Set("RepositoryRegistry", rr)
		c.Next()
	})
}

func (s *Server) SetLogger(logger *zap.Logger) {
	s.router.Use(func(c *gin.Context) {
		c.Set("Logger", logger)
		c.Next()
	})

	// s.router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	// s.router.Use(ginzap.RecoveryWithZap(logger, true))
}

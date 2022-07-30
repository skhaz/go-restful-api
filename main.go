package main

import (
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"skhaz.dev/rest/controller"
	"skhaz.dev/rest/database"
	"skhaz.dev/rest/repository"
)

func main() {
	viper.AutomaticEnv()

	logger, _ := zap.NewDevelopment()

	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		viper.Get("POSTGRES_HOST"),
		viper.Get("POSTGRES_PORT"),
		viper.Get("POSTGRES_USER"),
		viper.Get("POSTGRES_PASSWORD"),
		viper.Get("POSTGRES_DB"),
	)

	db, err := database.Connect(dsn)
	if err != nil {
		panic(err)
	}

	registry := repository.NewRepositoryRegistry(
		db,
		&repository.WorkspaceRepository{},
	)

	server := controller.InitServer()
	server.SetRepositoryRegistry(registry)
	server.SetLogger(logger)
	server.Run()
}

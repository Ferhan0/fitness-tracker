package main

import (
	"github.com/Ferhan0/fitness-tracker/initializers"
	"github.com/Ferhan0/fitness-tracker/routes"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncToDb()
}

func main() {
	router := gin.Default()
	routes.SetupRoutes(router)
	router.Run()
}

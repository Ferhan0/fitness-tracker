package routes

import (
	"github.com/Ferhan0/fitness-tracker/controllers"
	"github.com/Ferhan0/fitness-tracker/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.POST("/signup", controllers.SignUp)
	router.POST("/login", controllers.Login)
	router.POST("/refresh", controllers.RefreshToken)
	router.GET("/validate", middleware.RequireAuth, controllers.Validate)
}

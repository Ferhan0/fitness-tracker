package routes

import (
	"github.com/Ferhan0/fitness-tracker/controllers"
	"github.com/Ferhan0/fitness-tracker/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	//Auth(Korumasız) rotaları
	router.POST("/signup", controllers.SignUp)
	router.POST("/login", controllers.Login)
	router.POST("/refresh", controllers.RefreshToken)
	router.GET("/validate", middleware.RequireAuth, controllers.Validate)

	//Workout(Korumalı) rotaları
	workoutRoutes := router.Group("/workouts")
	workoutRoutes.Use(middleware.RequireAuth)
	{
		workoutRoutes.POST("/", controllers.CreateWorkout)
		workoutRoutes.GET("/", controllers.GetWorkouts)
		workoutRoutes.PUT("/:id", controllers.UpdateWorkout)
		workoutRoutes.DELETE("/:id", controllers.DeleteWorkout)
	}
}

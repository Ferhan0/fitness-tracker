package controllers

import (
	"github.com/Ferhan0/fitness-tracker/initializers"
	"github.com/Ferhan0/fitness-tracker/models"
	"github.com/gin-gonic/gin"
)

func CreateWorkout(c *gin.Context) {
	var body struct {
		Title string
		Date  string
		Notes string
	}
	if c.Bind(&body) != nil {
		c.JSON(400, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	user, _ := c.Get("user")
	currentUser := user.(models.User)

	workout := models.Workout{

		UserID: currentUser.ID,
		Title:  body.Title,
		Date:   body.Date,
		Notes:  body.Notes,
	}
	initializers.DB.Create(&workout)
	c.JSON(201, gin.H{"workout": workout})

}

func GetWorkouts(c *gin.Context) {

}

func UpdateWorkout(c *gin.Context) {

}

func DeleteWorkout(c *gin.Context) {

}

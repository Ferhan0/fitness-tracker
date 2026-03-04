package initializers

import "github.com/Ferhan0/fitness-tracker/models"

func SyncToDb() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Workout{})
}

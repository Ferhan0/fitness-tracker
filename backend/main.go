package main

import (
	"github.com/Ferhan0/fitness-tracker/initializers"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncToDb()
}

func main() {

}

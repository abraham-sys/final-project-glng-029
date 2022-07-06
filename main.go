package main

import (
	"final-project/database"
	"final-project/router"
	"os"

	"github.com/joho/godotenv"
)

var PORT = os.Getenv("PORT")

func main() {
	env := os.Getenv("FOO_ENV")
	if env == "" {
		env = "development"
	}

	godotenv.Load(".env." + env)

	database.StartDB()
	r := router.StartApp()
	r.Run(":3000")
}

package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/josephmakin/monerochat/api"
	"github.com/josephmakin/monerochat/handlers"
	v1 "github.com/josephmakin/monerochat/routes/v1"
	"github.com/josephmakin/monerohub/services"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Printf("Error loading .env file: %v", err)
    }

	api.Endpoint = os.Getenv("MONEROHUB_ENDPOINT")
	handlers.CallbackURL = os.Getenv("CALLBACK_URL")

	collections := map[string]string{
		"donations": "donations",
	}
    err = services.InitMongo(os.Getenv("MONGO_URI"), "monerochat", collections)
    if err != nil {
        log.Fatalf("Error connecting to mongo: %v", err)
    }

	router := gin.Default()
	v1.SetupRoutes(router)
	router.Run(":5000")
}

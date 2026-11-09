package main

import (
	"log"
	"tiket-kereta-be/config"
	"tiket-kereta-be/models"
	"tiket-kereta-be/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// load file ..env
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: File .env not found")
	}

	// connect to PostgreSQl db
	config.ConnectDB()

	// auto migrate tabel User (automatic create table in Postgres if not exist)
	config.DB.AutoMigrate(&models.User{}, &models.Station{})

	// init Gin Router
	r := gin.Default()
	r.Use(cors.Default())

	// route API
	routes.SetupRouter(r)

	// run server in port 3001
	r.Run(":3001")
}
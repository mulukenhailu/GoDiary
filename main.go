package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mulukenhailu/Diary_api/controller"
	"github.com/mulukenhailu/Diary_api/database"
	"github.com/mulukenhailu/Diary_api/middleware"
	"github.com/mulukenhailu/Diary_api/model"
)

func main() {
	loadEnv()
	loadDatabase()
	ServerApplication()
}

func loadDatabase() {
	database.Connect()
	database.Database.AutoMigrate(&model.User{})
	database.Database.AutoMigrate(&model.Entry{})
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func ServerApplication() {
	router := gin.Default()

	publicRoutes := router.Group("/auth")
	publicRoutes.POST("/register", controller.Register)
	publicRoutes.POST("/login", controller.Login)

	protectedRoutes := router.Group("/api")
	protectedRoutes.Use(middleware.JWTAuthMiddleware())
	protectedRoutes.GET("/entry", controller.GetAllEntry)
	protectedRoutes.POST("/entry", controller.AddEntry)

	router.Run(":8000")
	fmt.Println(" server is running on port 8080")
}

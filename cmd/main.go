package main

import (
    "pet-adoption-api/internal/config"
    "pet-adoption-api/internal/database"
    "pet-adoption-api/internal/handlers"
    "pet-adoption-api/internal/middleware"

    "github.com/gin-gonic/gin"
)

func main() {
    
    config.LoadEnv()

   
    database.Connect()

    
    r := gin.Default()

    auth := r.Group("/auth")
    {
        auth.POST("/register", handlers.Register)
        auth.POST("/login", handlers.Login)
    }

    petRoutes := r.Group("/pets")
    {
        petRoutes.GET("/", handlers.GetPets) 
        petRoutes.POST("/", middleware.AuthMiddleware(), middleware.AdminOnly(), handlers.CreatePet)
        petRoutes.DELETE("/:id", middleware.AuthMiddleware(), middleware.AdminOnly(), handlers.DeletePet)
    }

    shelterRoutes := r.Group("/shelters")
    {
        shelterRoutes.GET("/", handlers.GetShelters) 
        shelterRoutes.POST("/", middleware.AuthMiddleware(), middleware.AdminOnly(), handlers.CreateShelter)
    }

    r.Run(":8080")
}

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"pet-adoption-api/internal/config"
	"pet-adoption-api/internal/database"
	"pet-adoption-api/internal/handlers"
	"pet-adoption-api/internal/middleware"
	"pet-adoption-api/internal/worker"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load env vars (.env or system)
	config.LoadEnv()

	// Connect DB
	database.Connect()

	// Init JWT managers for handlers + middleware
	handlers.InitAuth()
	middleware.InitAuthMiddleware()

	// Create adoption worker with buffered channel
	aw := worker.NewAdoptionWorker(100)

	// Context for graceful shutdown (worker listens on this)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Start worker in background
	go aw.Start(ctx)

	// Gin router
	r := gin.Default()

	// Give worker channel to handlers (for pushing events)
	handlers.AdoptionEvents = aw.Events

	// Auth routes
	auth := r.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	// Pets routes
	petRoutes := r.Group("/pets")
	{
		petRoutes.GET("/", handlers.GetPets)
		petRoutes.POST("/", middleware.AuthMiddleware(), middleware.AdminOnly(), handlers.CreatePet)
		petRoutes.DELETE("/:id", middleware.AuthMiddleware(), middleware.AdminOnly(), handlers.DeletePet)
	}

	// Shelters routes
	shelterRoutes := r.Group("/shelters")
	{
		shelterRoutes.GET("/", handlers.GetShelters)
		shelterRoutes.POST("/", middleware.AuthMiddleware(), middleware.AdminOnly(), handlers.CreateShelter)
	}

	// Adoption routes (protected)
	adoptionRoutes := r.Group("/adoptions", middleware.AuthMiddleware())
	{
		// user creates a request
		adoptionRoutes.POST("/:petID/apply", handlers.ApplyForAdoption)

		// user sees only their own requests
		adoptionRoutes.GET("/my", handlers.GetMyAdoptions)

		// shelter owner/admin sees requests for their pets
		adoptionRoutes.GET("/shelter", middleware.ShelterOnly(), handlers.GetShelterAdoptions)

		// shelter owner/admin approve or reject
		adoptionRoutes.PATCH("/:id/approve", middleware.ShelterOnly(), handlers.ApproveAdoption)
		adoptionRoutes.PATCH("/:id/reject", middleware.ShelterOnly(), handlers.RejectAdoption)
	}

	// Read port from env, default 8080
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("server listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

	// When r.Run exits, context will be cancelled by signal, worker stops via ctx
}

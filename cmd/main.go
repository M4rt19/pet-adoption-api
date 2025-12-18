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
    "pet-adoption-api/internal/models"
	"pet-adoption-api/internal/worker"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load env vars (.env or system)
	config.LoadEnv()

	// Connect DB
	database.Connect()

	// Init JWT managers for handlers + middleware
	handlers.InitAuth()
	middleware.InitAuthMiddleware()

    // Temporary code to create users and shelters
	createTestData()

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
		petRoutes.GET("/:id", handlers.GetPetByID)
		petRoutes.POST("/", middleware.AuthMiddleware(), middleware.AdminOnly(), handlers.CreatePet)
		petRoutes.PUT("/:id", middleware.AuthMiddleware(), middleware.AdminOnly(), handlers.UpdatePet)
		petRoutes.DELETE("/:id", middleware.AuthMiddleware(), middleware.AdminOnly(), handlers.DeletePet)
	}

	// Shelters routes
	shelterRoutes := r.Group("/shelters")
	{
		shelterRoutes.GET("/", handlers.GetShelters)
		shelterRoutes.GET("/:id", handlers.GetShelterByID)
		shelterRoutes.POST("/", middleware.AuthMiddleware(), middleware.AdminOnly(), handlers.CreateShelter)
		shelterRoutes.PUT("/:id", middleware.AuthMiddleware(), middleware.AdminOnly(), handlers.UpdateShelter)
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

func createTestData() {
	users := []models.User{
		{Name: "Test User 1", Email: "user1@example.com", PasswordHash: "password", Role: "user"},
		{Name: "Test User 2", Email: "user2@example.com", PasswordHash: "password", Role: "user"},
		{Name: "Shelter Owner 1", Email: "shelter1@example.com", PasswordHash: "password", Role: "shelter"},
		{Name: "Shelter Owner 2", Email: "shelter2@example.com", PasswordHash: "password", Role: "shelter"},
	}

	for _, user := range users {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("failed to hash password: %v", err)
		}
		user.PasswordHash = string(hash)
		if err := database.DB.Create(&user).Error; err != nil {
			log.Printf("could not create user %s: %v", user.Name, err)
		}
	}

	var shelterOwner1 models.User
	database.DB.Where("email = ?", "shelter1@example.com").First(&shelterOwner1)

	var shelterOwner2 models.User
	database.DB.Where("email = ?", "shelter2@example.com").First(&shelterOwner2)


	shelters := []models.Shelter{
		{Name: "Happy Paws Shelter", Address: "123 Main St", Phone: "555-1234", OwnerUserID: shelterOwner1.ID},
		{Name: "Furry Friends Rescue", Address: "456 Oak Ave", Phone: "555-5678", OwnerUserID: shelterOwner2.ID},
	}

	for _, shelter := range shelters {
		if err := database.DB.Create(&shelter).Error; err != nil {
			log.Printf("could not create shelter %s: %v", shelter.Name, err)
		}
	}
}

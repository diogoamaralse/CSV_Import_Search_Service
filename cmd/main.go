package main

import (
	"ImportAndSearchCsvFile/internal/controller/users"
	"ImportAndSearchCsvFile/internal/service"

	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	userService := service.NewUserStore()
	router := gin.Default()
	users.NewUsersHandler(router, userService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

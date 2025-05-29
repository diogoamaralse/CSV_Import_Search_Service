package users

import (
	"ImportAndSearchCsvFile/internal/service"

	"github.com/gin-gonic/gin"
)

func (h *handler) GetUserHandler(store service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.Query("email")
		if email == "" {
			c.JSON(400, gin.H{"error": "Email parameter is required"})
			return
		}

		user, exists := store.GetUserByEmail(email)
		if !exists {
			c.JSON(404, gin.H{"error": "User not found"})
			return
		}

		c.JSON(200, user)
	}
}

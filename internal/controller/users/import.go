package users

import (
	"ImportAndSearchCsvFile/internal/service"

	"github.com/gin-gonic/gin"
)

func (h *handler) ImportUsersHandler(store service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(400, gin.H{"error": "Failed to get file from request"})
			return
		}

		f, err := file.Open()
		if err != nil {
			c.JSON(400, gin.H{"error": "Failed to open uploaded file"})
			return
		}
		defer f.Close()

		if err := store.ImportUsers(f); err != nil {
			c.JSON(400, gin.H{"error": "Failed to import users", "details": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Users imported successfully"})
	}
}

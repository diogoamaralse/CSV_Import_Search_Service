package users

import (
	"ImportAndSearchCsvFile/internal/service"

	"github.com/gin-gonic/gin"
)

type handler struct {
	UserStoreService service.Service
}

func NewUsersHandler(g *gin.Engine, userStoreService service.Service) {

	h := &handler{
		UserStoreService: userStoreService,
	}

	api := g.Group("/api/v1/user")
	{
		api.GET("", h.GetUserHandler(h.UserStoreService))
		api.POST("", h.ImportUsersHandler(h.UserStoreService))
	}
}

package users

import (
	"ImportAndSearchCsvFile/internal/service"
	"ImportAndSearchCsvFile/pkg/models"

	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter(h *handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	api := router.Group("/api/v1/user")
	{
		api.GET("", h.GetUserHandler(h.UserStoreService))
		api.POST("", h.ImportUsersHandler(h.UserStoreService))
	}
	return router
}

func TestUserHandler(t *testing.T) {
	mockService := &service.MockService{
		GetUserByEmailFn: func(email string) (models.User, bool) {
			if email == "test@example.com" {
				return models.User{
					ID:        1,
					FirstName: "Test",
					LastName:  "User",
					Email:     "test@example.com",
				}, true
			}
			return models.User{}, false
		},
	}

	h := &handler{UserStoreService: mockService}
	router := setupTestRouter(h)

	tests := []struct {
		name       string
		query      string
		wantStatus int
	}{
		{"existing user", "email=test@example.com", http.StatusOK},
		{"missing email", "", http.StatusBadRequest},
		{"nonexistent user", "email=nonexistent@example.com", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/v1/user?"+tt.query, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestImportHandler(t *testing.T) {
	mockService := &service.MockService{
		ImportUsersFn: func(reader io.Reader) error {
			return nil
		},
	}

	h := &handler{UserStoreService: mockService}
	router := setupTestRouter(h)

	t.Run("successful import", func(t *testing.T) {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "test.csv")
		part.Write([]byte("id,first_name,last_name,email\n1,Test,User,test@example.com"))
		writer.Close()

		req, _ := http.NewRequest("POST", "/api/v1/user", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("no file provided", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/v1/user", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

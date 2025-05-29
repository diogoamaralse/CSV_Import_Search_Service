package main

import (
	"ImportAndSearchCsvFile/internal/controller/users"
	"ImportAndSearchCsvFile/internal/service"
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Setenv("GIN_MODE", "test")
	code := m.Run()
	os.Exit(code)
}

func TestAPIEndpoints(t *testing.T) {
	userService := service.NewUserStore()
	router := setupRouter(userService)
	ts := httptest.NewServer(router)
	defer ts.Close()

	t.Run("test import and search", func(t *testing.T) {
		csvData := `id,first_name,last_name,email,created_at,deleted_at,merged_at,parent_user_id
1,Test,User,test@example.com,123456789,0,0,0`

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("file", "test.csv")
		assert.NoError(t, err)

		_, err = io.WriteString(part, csvData)
		assert.NoError(t, err)

		err = writer.Close()
		assert.NoError(t, err)

		req, err := http.NewRequest("POST", ts.URL+"/api/v1/user", body)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		resp.Body.Close()

		resp, err = http.Get(ts.URL + "/api/v1/user?email=test@example.com")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		defer resp.Body.Close()
		bodyBytes, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Contains(t, string(bodyBytes), `"email":"test@example.com"`)
	})
}

func setupRouter(userService service.Service) *gin.Engine {
	router := gin.Default()
	users.NewUsersHandler(router, userService)
	return router
}

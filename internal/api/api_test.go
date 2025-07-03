package api_test

import (
	"net/http"
	"net/http/httptest"
	"shop-dashboard/internal/api"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	router := api.NewRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "OK", w.Body.String())
}

func TestExample(t *testing.T) {
	router := api.NewRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/example", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{\"message\": \"Hello, World!\"}", w.Body.String())
}

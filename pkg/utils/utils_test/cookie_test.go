package utils_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"wtm-backend/pkg/utils"
)

func TestSetAndClearRefreshCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.GET("/set", func(c *gin.Context) {
		utils.SetRefreshCookie(c, "token123", "localhost", 3600, false)
	})

	router.GET("/clear", func(c *gin.Context) {
		utils.ClearRefreshCookie(c, "localhost", false)
	})

	t.Run("set cookie", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/set", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		cookies := w.Result().Cookies()
		assert.Len(t, cookies, 1)
		assert.Equal(t, "refresh_token", cookies[0].Name)
		assert.Equal(t, "token123", cookies[0].Value)
	})

	t.Run("clear cookie", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/clear", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		cookies := w.Result().Cookies()
		assert.Len(t, cookies, 1)
		assert.Equal(t, "refresh_token", cookies[0].Name)
		assert.Equal(t, "", cookies[0].Value)
		assert.Equal(t, -1, cookies[0].MaxAge)
	})
}

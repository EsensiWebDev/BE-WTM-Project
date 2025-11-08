package utils

import "github.com/gin-gonic/gin"

func SetRefreshCookie(c *gin.Context, refreshToken, domain string, maxAge int, secure bool) {
	c.SetCookie("refresh_token", refreshToken, maxAge, "/", domain, secure, true)
}

func ClearRefreshCookie(c *gin.Context, domain string, secure bool) {
	c.SetCookie("refresh_token", "", -1, "/", domain, secure, true)
}

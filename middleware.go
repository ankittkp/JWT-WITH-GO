package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Middleware : two routes that require authentication: /login and /logout as anybody can access that, Middleware will secure these routes
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := CheckValidToken(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		c.Next()
	}
}

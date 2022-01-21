package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Logout : Extract the JWT metadata. If true then delete the metadata, and so JWT invalid immediately.
func Logout(c *gin.Context) {
	au, err := GetTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	deleted, delErr := DeleteMetadataFromRedis(au.AccessUuid)
	if delErr != nil || deleted == 0 { //if any goes wrong
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	c.JSON(http.StatusOK, "Successfully logged out")
}

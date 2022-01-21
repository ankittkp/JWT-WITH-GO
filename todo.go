package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateTodo(c *gin.Context) {
	var td *Todo
	if err := c.ShouldBindJSON(&td); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "invalid json")
		return
	}
	tokenAuth, err := GetTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	userId, err := SearchMetadataInRedis(tokenAuth)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	td.UserID = userId

	//you can proceed to save the Todo to a database
	//but we will just return it to the caller here:
	c.JSON(http.StatusCreated, td)
}
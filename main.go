package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
)

var (
	router = gin.Default()
	client *redis.Client
	user   = User{
		ID:       1,
		Username: "username",
		Password: "password",
	}
	td = &TokenDetails{}
)

/*
Each time we run main.go redis will automatically connect
*/
func init() {
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	client = redis.NewClient(&redis.Options{
		Addr: dsn,
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("JWT with GO")
	router.POST("/api/v1/login", Login)
	router.POST("/api/v2/login", NewLogin)
	router.POST("/api/v2/todo", Middleware(), CreateTodo)
	router.POST("/api/v2/logout", Middleware(), Logout)
	router.POST("/api/v2/token/refresh", TokenRefresh)
	log.Fatal(router.Run(":8080"))
}

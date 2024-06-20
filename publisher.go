package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var (
	ctx = context.Background()
)

func main() {
	r := gin.Default()

	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})

	r.POST("/publish/:message", func(c *gin.Context) {
		message := c.Param("message")
		err := rdb.Publish(ctx, "channel", message).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Message published"})
	})

	r.Run(":8000")
}

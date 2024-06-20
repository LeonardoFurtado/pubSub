package main

import (
	"context"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var (
	cts = context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})
	messages []string
	mu       sync.Mutex
)

func main() {
	r := gin.Default()

	pubsub := rdb.Subscribe(cts, "channel")
	ch := pubsub.Channel()

	go func() {
		for msg := range ch {
			mu.Lock()
			messages = append(messages, msg.Payload)
			mu.Unlock()
		}
	}()

	r.GET("/messages", func(c *gin.Context) {
		mu.Lock()
		defer mu.Unlock()
		c.JSON(http.StatusOK, gin.H{"messages": messages})
	})

	r.Run(":8001")
}

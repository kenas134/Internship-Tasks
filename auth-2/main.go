package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Authentication API",
            "version": "1.0.0",
        })
    })
	router.POST("/register", registerHandler)

	if err := router.Run(":8080"); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
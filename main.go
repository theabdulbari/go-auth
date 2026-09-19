package main

import (
	"go-auth/database"
	"go-auth/handlers"
	"go-auth/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()

	r := gin.Default()

	auth := r.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	protected := r.Group("/api")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/profile", func(c *gin.Context) {
			userID := c.GetUint("user_id")
			c.JSON(http.StatusOK, gin.H{"user_id": userID})
		})
	}

	r.Run(":8080")
}

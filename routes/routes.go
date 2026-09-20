package routes

import (
	"time"

	"go-auth/handlers"
	"go-auth/middleware"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Rate limiter: 5 login attempts per minute, burst of 5
	loginLimiter := middleware.NewRateLimiter(rate.Every(12*time.Second), 5)
	refreshLimiter := middleware.NewRateLimiter(rate.Every(6*time.Second), 10)

	// Public routes
	auth := r.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", loginLimiter.Middleware(), handlers.Login) // rate-limited
		auth.POST("/refresh", refreshLimiter.Middleware(), handlers.Refresh)
		auth.POST("/logout", handlers.Logout)
	}

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.AuthRequired())
	{
		api.GET("/profile", handlers.Profile)
		api.GET("/users", handlers.ListUsers)
		api.POST("/logout-all", handlers.LogoutAll)
	}

	return r
}
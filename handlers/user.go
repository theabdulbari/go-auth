package handlers

import (
	"net/http"
	"strconv"

	"go-auth/database"
	"go-auth/models"

	"github.com/gin-gonic/gin"
)

// Profile returns the currently authenticated user's info.
// Requires: AuthRequired() middleware (sets "user_id" in context).
func Profile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"joined":   user.CreatedAt,
	})
}

// ListUsers returns all users (without password fields).
// Requires: AuthRequired() middleware.

func ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	var users []models.User
	var total int64

	database.DB.Model(&models.User{}).Count(&total)
	if err := database.DB.
		Order("id ASC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// Build a lightweight response slice — never leak password hashes
	type UserResponse struct {
		ID        uint   `json:"id"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		CreatedAt string `json:"created_at"`
	}

	resp := make([]UserResponse, 0, len(users))
	for _, u := range users {
		resp = append(resp, UserResponse{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"page":  page,
		"limit": limit,
		"total": total,
		"users": resp,
	})
}
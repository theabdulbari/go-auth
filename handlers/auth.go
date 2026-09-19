package handlers

import (
	"net/http"
	"time"

	"go-auth/database"
	"go-auth/models"
	"go-auth/utils"

	"github.com/gin-gonic/gin"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// helper: issue both tokens and persist refresh token
func issueTokenPair(c *gin.Context, user models.User) (string, string, error) {
	access, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		return "", "", err
	}

	refresh, err := utils.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	rt := models.RefreshToken{
		UserID:    user.ID,
		Token:     refresh,
		ExpiresAt: time.Now().Add(utils.RefreshTokenExpiry()),
		UserAgent: c.GetHeader("User-Agent"),
		IP:        c.ClientIP(),
	}
	if err := database.DB.Create(&rt).Error; err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashed, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: hashed,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	access, refresh, err := issueTokenPair(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not issue tokens"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":       "User registered",
		"user_id":       user.ID,
		"access_token":  access,
		"refresh_token": refresh,
		"expires_in":    int(utils.RefreshTokenExpiry().Seconds()),
	})
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !utils.CheckPassword(user.Password, input.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	access, refresh, err := issueTokenPair(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not issue tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
		"user":          user.Username,
		"expires_in":    int(utils.RefreshTokenExpiry().Seconds()),
	})
}

// Refresh rotates the refresh token (old one is revoked)
func Refresh(c *gin.Context) {
	var input RefreshInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var rt models.RefreshToken
	err := database.DB.
		Where("token = ? AND revoked = ?", input.RefreshToken, false).
		First(&rt).Error
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	if time.Now().After(rt.ExpiresAt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired"})
		return
	}

	// Revoke old token (rotation)
	rt.Revoked = true
	database.DB.Save(&rt)

	// Fetch user
	var user models.User
	if err := database.DB.First(&user, rt.UserID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	access, refresh, err := issueTokenPair(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not issue tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

// Logout revokes a refresh token
func Logout(c *gin.Context) {
	var input RefreshInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Model(&models.RefreshToken{}).
		Where("token = ?", input.RefreshToken).
		Update("revoked", true)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}
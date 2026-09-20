package handlers

import (
	"net/http"
	"time"
	"errors"

	"go-auth/database"
	"go-auth/models"
	"go-auth/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

// helper: issue both tokens and persist refresh token (not used after added issueTokenPairTx). you may delete this (primary integration)
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

// issueTokenPairTx issues both tokens and persists the refresh token using the given tx.
func issueTokenPairTx(tx *gorm.DB, c *gin.Context, user models.User) (string, string, error) {
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
	if err := tx.Create(&rt).Error; err != nil {
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

	// Begin transaction: user + refresh token are created atomically
	tx := database.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// Rollback on any failure; commit on success
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()


	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}
	access, refresh, err := issueTokenPairTx(tx, c, user)
	if err != nil {
		tx.Rollback() // undoes the user creation
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Could not issue tokens",
			"detail": err.Error(), // remove in production
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit"})
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

	tx := database.DB.Begin()

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !utils.CheckPassword(user.Password, input.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	access, refresh, err := issueTokenPairTx(tx, c, user)
	if err != nil {
		tx.Rollback()
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

	// Look up the token (revoked or not, we need to know)
	var existing models.RefreshToken
	err := database.DB.
		Where("token = ?", input.RefreshToken).
		First(&existing).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Reuse detection: a revoked token was presented.
	// This means either the token was stolen OR the user's token was already rotated.
	// Safest action: revoke ALL of this user's refresh tokens -> force re-login.
	if existing.Revoked {
		if err := database.DB.
			Model(&models.RefreshToken{}).
			Where("user_id = ? AND revoked = ?", existing.UserID, false).
			Update("revoked", true).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke tokens"})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Refresh token reuse detected. All sessions revoked.",
		})
		return
	}

	// Expired?
	if time.Now().After(existing.ExpiresAt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired"})
		return
	}

	// Load the user
	var user models.User
	if err := database.DB.First(&user, existing.UserID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Atomic rotation: revoke old + persist new in one transaction
	tx := database.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Revoke the old token inside the tx
	if err := tx.Model(&models.RefreshToken{}).
		Where("id = ?", existing.ID).
		Update("revoked", true).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke old token"})
		return
	}

	// Issue a fresh pair (write new refresh token via tx)
	access, refresh, err := issueTokenPairTx(tx, c, user)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Could not issue tokens",
			"detail": err.Error(), // remove in production
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit"})
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

// LogoutAll revokes every active refresh token for the authenticated user.
// Mount it on a protected route so you know who the user is.
func LogoutAll(c *gin.Context) {
	userID := c.GetUint("user_id")

	if err := database.DB.
		Model(&models.RefreshToken{}).
		Where("user_id = ? AND revoked = ?", userID, false).
		Update("revoked", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke sessions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All sessions revoked"})
}
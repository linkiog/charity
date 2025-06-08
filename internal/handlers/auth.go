package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/linkiog/charity/internal/config"
	"github.com/linkiog/charity/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB     *gorm.DB
	Config *config.Config
}

func NewAUthHAndler(db *gorm.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{DB: db, Config: cfg}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := models.User{Username: req.Username, Email: req.Email, Password: string(hash)}

	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "user exists"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "registered"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "bad request"})
		return
	}
	var user models.User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	})
	tokStr, _ := token.SignedString([]byte(h.Config.JWTSecret))

	c.JSON(http.StatusOK, gin.H{"token": tokStr})
}

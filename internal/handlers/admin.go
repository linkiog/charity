package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/linkiog/charity/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminHandler struct {
	DB *gorm.DB
}

func NewAdminHandler(db *gorm.DB) *AdminHandler {
	return &AdminHandler{DB: db}
}

func (h *AdminHandler) CreateAdmin(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"    binding:"required,email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err := h.DB.Where("email = ?", req.Email).First(&user).Error

	if err == nil {
		if user.Role == models.RoleAdmin {
			c.JSON(http.StatusConflict, gin.H{"error": "already admin"})
			return
		}
		user.Role = models.RoleAdmin
		if err := h.DB.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot promote"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"admin_id": user.ID, "promoted": true})
		return
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		if req.Username == "" || req.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "username and password required for new admin",
			})
			return
		}

		hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		user = models.User{
			Username: req.Username,
			Email:    req.Email,
			Password: string(hash),
			Role:     models.RoleAdmin,
		}
		if err := h.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"admin_id": user.ID, "created": true})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
}

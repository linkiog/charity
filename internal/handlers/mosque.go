package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/linkiog/charity/internal/dto"
	"github.com/linkiog/charity/models"
	"gorm.io/gorm"
)

type MosqueHandler struct {
	DB *gorm.DB
}

func NewMosqueHandler(db *gorm.DB) *MosqueHandler {
	return &MosqueHandler{DB: db}
}

func (h *MosqueHandler) Create(c *gin.Context) {
	var req struct {
		Name       string `json:"name" binding:"required"`
		City       string `json:"city"`
		Region     string `json:"region"`
		Requisites string `json:"requisites"`
		AdminID    uint   `json:"admin_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mosque := models.Mosque{
		Name:       req.Name,
		City:       req.City,
		Region:     req.Region,
		Requisites: req.Requisites,
		AdminID:    req.AdminID,
	}
	if err := h.DB.Create(&mosque).Error; err != nil {
		log.Printf("failed to create mosque: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create mosque"})
		return
	}

	resp := dto.MosqueResponse{
		ID:         mosque.ID,
		Name:       mosque.Name,
		City:       mosque.City,
		Region:     mosque.Region,
		Requisites: mosque.Requisites,
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *MosqueHandler) List(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetUint("user_id")

	type row struct {
		ID        uint
		Name      string
		City      string
		Region    string
		Need      float64
		Collected float64
	}

	dbq := h.DB.
		Table("mosques").
		Select(`
			mosques.id,
			mosques.name,
			mosques.city,
			mosques.region,
			COALESCE(SUM(products.price * products.need), 0)      AS need,
			COALESCE(SUM(products.price * products.purchased), 0) AS collected`,
		).
		Joins("LEFT JOIN products ON products.mosque_id = mosques.id").
		Group("mosques.id")

	if role == models.RoleAdmin {
		dbq = dbq.Where("mosques.admin_id = ?", userID)
	}

	var rows []row
	if err := dbq.Scan(&rows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	resp := make([]dto.MosqueListItemWithProgress, 0, len(rows))
	for _, r := range rows {
		resp = append(resp, dto.MosqueListItemWithProgress{
			ID:        r.ID,
			Name:      r.Name,
			City:      r.City,
			Region:    r.Region,
			Need:      r.Need,
			Collected: r.Collected,
			Remaining: r.Need - r.Collected,
		})
	}

	c.JSON(http.StatusOK, resp)
}

func (h *MosqueHandler) GetWithProducts(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("mosqueID"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mosque id"})
		return
	}

	var mosque models.Mosque
	db := h.DB.Preload("Products")
	role := c.GetString("role")
	if role == models.RoleAdmin {
		userID := c.GetUint("user_id")
		db = db.Where("admin_id = ?", userID)
	}

	if err := db.First(&mosque, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "mosque not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	resp := dto.MosqueFull{
		ID:         mosque.ID,
		Name:       mosque.Name,
		City:       mosque.City,
		Region:     mosque.Region,
		Requisites: mosque.Requisites,
		Products:   make([]dto.ProductItem, len(mosque.Products)),
	}
	for i, p := range mosque.Products {
		resp.Products[i] = dto.ProductItem{
			ID:        p.ID,
			Name:      p.Name,
			Price:     p.Price,
			Need:      p.Need,
			Purchased: p.Purchased,
		}
	}

	c.JSON(http.StatusOK, resp)
}

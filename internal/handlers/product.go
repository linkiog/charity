package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/linkiog/charity/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductHandler struct {
	DB *gorm.DB
}

func NewProductHandler(db *gorm.DB) *ProductHandler {
	return &ProductHandler{DB: db}
}

// func (h *ProductHandler) ListByMosque(c *gin.Context) {
// 	mid, err := strconv.ParseUint(c.Param("mosqueID"), 10, 64)
// 	if err != nil || mid == 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mosque_id"})
// 		return
// 	}

// 	if c.GetString("role") == models.RoleAdmin && c.GetUint("mosque_id") != uint(mid) {
// 		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
// 		return
// 	}

// 	var products []models.Product
// 	h.DB.Where("mosque_id = ?", mid).Find(&products)

// 	c.JSON(http.StatusOK, products)
// }

func (h *ProductHandler) Buy(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("productID"))

	var req struct {
		Qty int `json:"qty" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var prod models.Product
	if err := h.DB.Clauses(clause.Locking{Strength: "UPDATE"}).First(&prod, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	if prod.Purchased+req.Qty > prod.Need {
		c.JSON(http.StatusBadRequest, gin.H{"error": "exceeds required amount"})
		return
	}
	uid := c.GetUint("user_id")

	don := models.Donation{
		ProductID: uint(id),
		Qty:       req.Qty,
		Amount:    prod.Price * float64(req.Qty),
	}
	if uid != 0 {
		don.UserID = uid
	}
	if err := h.DB.Create(&don).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot save donation"})
		return
	}

	if err := h.DB.Model(&prod).UpdateColumn("purchased", gorm.Expr("purchased + ?", req.Qty)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot update purchased"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *ProductHandler) CreateForMosque(c *gin.Context) {
	mid, err := strconv.ParseUint(c.Param("mosqueID"), 10, 64)
	if err != nil || mid == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mosque_id"})
		return
	}

	if c.GetString("role") == models.RoleAdmin && c.GetUint("mosque_id") != uint(mid) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var mosque models.Mosque
	if err := h.DB.First(&mosque, mid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "mosque not found"})
		return
	}

	var req struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description"`
		Price       float64 `json:"price" binding:"required"`
		Need        int     `json:"need" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := models.Product{
		MosqueID:    uint(mid),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Need:        req.Need,
	}

	if err := h.DB.Create(&product).Error; err != nil {
		log.Printf("Failed to create products: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create product"})
		return
	}

	c.JSON(http.StatusCreated, product)
}

package handlers

import (
	"fmt"
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

	mosqueID, err := strconv.ParseUint(c.Param("mosqueID"), 10, 64)
	if err != nil || mosqueID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mosque id"})
		return
	}
	productID, err := strconv.ParseUint(c.Param("productID"), 10, 64)
	if err != nil || productID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	var req struct {
		Qty int `json:"qty" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uid := c.GetUint("user_id")

	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		var prod models.Product

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND mosque_id = ?", productID, mosqueID).
			First(&prod).Error; err != nil {
			return fmt.Errorf("product not found")
		}

		if prod.Purchased+req.Qty > prod.Need {
			return fmt.Errorf("exceeds required amount")
		}

		don := models.Donation{
			ProductID: uint(productID),
			Qty:       req.Qty,
			Amount:    prod.Price * float64(req.Qty),
		}
		if uid != 0 {
			don.UserID = uid
		}
		if err := tx.Create(&don).Error; err != nil {
			fmt.Printf("cannot save donation %v", err)
			return fmt.Errorf("cannot save donation")
		}

		if err := tx.Model(&prod).
			UpdateColumn("purchased", gorm.Expr("purchased + ?", req.Qty)).
			Error; err != nil {
			return fmt.Errorf("cannot update purchased")
		}

		return nil
	}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

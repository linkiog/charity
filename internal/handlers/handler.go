package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/linkiog/charity/internal/config"
	"github.com/linkiog/charity/internal/middleware"
	"github.com/linkiog/charity/models"
	"gorm.io/gorm"
)

func Handler(cfg *config.Config, gorm *gorm.DB, r *gin.Engine) {

	authH := NewAuthHandler(gorm, cfg)
	adminH := NewAdminHandler(gorm)
	mosqH := NewMosqueHandler(gorm)
	prodH := NewProductHandler(gorm)

	api := r.Group("/api")

	authG := api.Group("/auth")
	{
		authG.POST("/register", authH.Register)
		authG.POST("/login", authH.Login)
	}

	auth := api.Group("", middleware.JWTAuth(cfg.JWTSecret))

	mosqSA := auth.Group("/mosques", middleware.RoleAuth(models.RoleSuperAdmin))
	{
		mosqSA.POST("", mosqH.Create)
		// mosqSA.PUT("/:id", mosqH.Update)
		// mosqSA.PUT("/:id	/admin", mosqH.SetAdmin)
		// mosqSA.DELETE("/:id", mosqH.Delete)
	}

	auth.POST("/admin/create", middleware.RoleAuth(models.RoleSuperAdmin), adminH.CreateAdmin)

	auth.GET("/mosques", middleware.RoleAuth(models.RoleSuperAdmin, models.RoleAdmin, models.RoleUser), mosqH.List)

	mosq := auth.Group("/mosques/:mosqueID", middleware.RoleAuth(models.RoleSuperAdmin, models.RoleAdmin))
	{
		mosq.GET("", mosqH.GetWithProducts)
		//mosq.GET("/products", prodH.ListByMosque)
		mosq.POST("/products", prodH.CreateForMosque)
		mosq.POST("/:productID/buy", prodH.Buy)
	}

}

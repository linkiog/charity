package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/linkiog/charity/internal/config"
	"gorm.io/gorm"
)

func Handler(cfg *config.Config, gorm *gorm.DB, r *gin.Engine) {
	authH := NewAUthHAndler(gorm, cfg)

	api := r.Group("/api")
	{
		api.POST("/register", authH.Register)
		api.POST("/login", authH.Login)
	}

}

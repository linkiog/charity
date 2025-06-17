package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/linkiog/charity/internal/config"
	"github.com/linkiog/charity/internal/db"
	"github.com/linkiog/charity/internal/handlers"
)

func main() {
	cfg := config.Load()
	if cfg.DBUrl == "" || cfg.JWTSecret == "" {
		log.Fatal("Variables bot sets")
	}
	gormDB := db.New(cfg)
	r := gin.Default()
	handlers.Handler(cfg, gormDB, r)
	r.Run(":8081")

}

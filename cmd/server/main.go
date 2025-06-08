package main

import (
	"fmt"
	"log"

	"github.com/linkiog/charity/internal/config"
	"github.com/linkiog/charity/internal/db"
)

func main() {

	cfg := config.Load()
	if cfg.DBUrl == "" || cfg.JWTSecret == "" {
		log.Fatal("Variables bot sets")
	}

	gormDB := db.New(cfg)
	fmt.Println(gormDB)

}

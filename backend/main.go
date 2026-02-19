package main

import (
	"log"

	"red-packet/config"
	"red-packet/database"
	"red-packet/router"
	"red-packet/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := database.Init(cfg); err != nil {
		log.Fatalf("failed to init database: %v", err)
	}
	log.Println("database connected and migrated")

	service.InitUserService(cfg.JWT.Secret, cfg.JWT.ExpireHours)

	r := router.NewRouter()
	r.Run(":" + cfg.Server.Port)
}

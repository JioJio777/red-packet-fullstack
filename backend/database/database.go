package database

import (
	"red-packet/config"
	"red-packet/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(cfg *config.Config) error {
	db, err := gorm.Open(mysql.Open(cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(
		&model.User{},
		&model.RedPacket{},
		&model.RedPacketRecord{},
		&model.Transaction{},
	)
	if err != nil {
		return err
	}

	DB = db
	return nil
}

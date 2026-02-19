package repository

import (
	"red-packet/database"
	"red-packet/model"
)

func CreateUser(user *model.User) error {
	return database.DB.Create(user).Error
}

func GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := database.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByID(id uint64) (*model.User, error) {
	var user model.User
	err := database.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

package main

import "gorm.io/gorm"

type User struct {
	ID    uint
	Name  string
	Email string
}

func GetUsers(db *gorm.DB, name string) ([]User, error) {
	var users []User
	db.Where("name = ?", name).Find(&users)
	return users, nil
}

package main

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"column:email;type:varchar;size:255"`
	Name     string `gorm:"column:name;type:varchar;size:255"`
	Password string `gorm:"column:password;type:varchar;size:255"`
}

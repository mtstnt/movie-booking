package main

import "gorm.io/gorm"

type Employee struct {
	gorm.Model
	Username string `gorm:"column:username;type:varchar;size:255"`
	Password string `gorm:"column:password;type:varchar;size:255"`
}

package main

import (
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{db}
}

func (r Repository) GetUsers() ([]User, error) {
	var users []User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r Repository) GetUser(id uint32) (User, error) {
	var user User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (r Repository) GetUserByEmail(email string) (User, error) {
	var user User
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		return User{}, err
	}
	return user, nil
}

func (r Repository) CreateUser(user *User) error {
	return r.db.Create(user).Error
}

func (r Repository) UpdateUser(user *User) error {
	return r.db.Updates(user).Error
}

func (r Repository) DeleteUser(user *User) error {
	user.DeletedAt = gorm.DeletedAt{
		Time:  time.Now(),
		Valid: true,
	}
	return r.db.Updates(user).Error
}

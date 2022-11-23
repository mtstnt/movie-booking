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

func (r Repository) GetEmployees() ([]Employee, error) {
	var employee []Employee
	if err := r.db.Find(&employee).Error; err != nil {
		return nil, err
	}
	return employee, nil
}

func (r Repository) GetEmployee(id uint32) (Employee, error) {
	var employee Employee
	if err := r.db.First(&employee, "id = ?", id).Error; err != nil {
		return employee, err
	}
	return employee, nil
}

func (r Repository) GetEmployeeByUsername(username string) (Employee, error) {
	var employee Employee
	if err := r.db.First(&employee, "username = ?", username).Error; err != nil {
		return Employee{}, err
	}
	return employee, nil
}

func (r Repository) CreateEmployee(employee *Employee) error {
	return r.db.Create(employee).Error
}

func (r Repository) UpdateEmployee(employee *Employee) error {
	return r.db.Updates(employee).Error
}

func (r Repository) DeleteEmployee(employee *Employee) error {
	employee.DeletedAt = gorm.DeletedAt{
		Time:  time.Now(),
		Valid: true,
	}
	return r.db.Updates(employee).Error
}

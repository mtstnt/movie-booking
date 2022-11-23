package main

import "golang.org/x/crypto/bcrypt"

type Service struct {
	employeeRepository Repository
}

func NewService(employeeRepository Repository) Service {
	return Service{employeeRepository}
}

func (s Service) AuthenticateEmployee(email string, rawPassword string) (Employee, error) {
	employee, err := s.employeeRepository.GetEmployeeByUsername(email)
	if err != nil {
		// TODO: Convert into prepared error to display in gateway.
		return Employee{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(employee.Password), []byte(rawPassword)); err != nil {
		// TODO: Convert into prepared error to display in gateway.
		return Employee{}, err
	}

	return employee, nil
}

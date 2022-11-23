package main

import "golang.org/x/crypto/bcrypt"

type Service struct {
	userRepository Repository
}

func NewService(userRepository Repository) Service {
	return Service{userRepository}
}

func (s Service) AuthenticateEmployee(email string, rawPassword string) (User, error) {
	user, err := s.userRepository.GetUserByEmail(email)
	if err != nil {
		// TODO: Convert into prepared error to display in gateway.
		return User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(rawPassword)); err != nil {
		// TODO: Convert into prepared error to display in gateway.
		return User{}, err
	}

	return user, nil
}

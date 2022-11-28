package main

type Service struct {
	movieRepository Repository
}

func NewService(movieRepository Repository) Service {
	return Service{movieRepository}
}

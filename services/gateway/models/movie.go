package models

import "time"

type Director struct {
	ID   uint32
	Name string
}

type Actor struct {
	ID   uint32
	Name string
}

type Movie struct {
	ID          uint32
	Title       string
	Synopsis    string
	ReleaseDate time.Time
	Director    Director
}

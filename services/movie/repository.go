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

func (r Repository) GetMovies() ([]Movie, error) {
	var movie []Movie
	if err := r.db.
		Preload("Casts").
		Preload("Director").
		Find(&movie).Error; err != nil {
		return nil, err
	}
	return movie, nil
}

func (r Repository) GetMovie(id uint32) (Movie, error) {
	var movie Movie
	if err := r.db.
		Preload("Casts").
		Preload("Director").
		First(&movie, "id = ?", id).Error; err != nil {
		return movie, err
	}
	return movie, nil
}

func (r Repository) CreateMovie(movie *Movie) error {
	return r.db.Create(movie).Error
}

func (r Repository) UpdateMovie(movie *Movie) error {
	return r.db.Updates(movie).Error
}

func (r Repository) DeleteMovie(movie *Movie) error {
	movie.DeletedAt = gorm.DeletedAt{
		Time:  time.Now(),
		Valid: true,
	}
	return r.db.Updates(movie).Error
}

package main

import (
	"movie/movie/pb"

	"gorm.io/gorm"
)

type Director struct {
	gorm.Model
	Name string `gorm:"column:name;type:varchar;size:255"`
}

type Actor Director

type Movie struct {
	gorm.Model
	Title       string `gorm:"column:title;type:varchar;size:255"`
	Synopsis    string `gorm:"column:synopsis;type:text;"`
	Director    Director
	DirectorID  uint
	Casts       []Actor `gorm:"many2many:movie_casts"`
	ReleaseDate uint64
}

func MovieToProto(movie *Movie) *pb.Movie {
	return &pb.Movie{
		Id:          uint32(movie.ID),
		Title:       movie.Title,
		Synopsis:    movie.Synopsis,
		ReleaseDate: movie.ReleaseDate,
		Director: &pb.Director{
			Id:   uint32(movie.DirectorID),
			Name: movie.Director.Name,
		},
		Casts: func() []*pb.Actor {
			protoCasts := make([]*pb.Actor, 0)
			for _, actor := range movie.Casts {
				protoCasts = append(protoCasts, &pb.Actor{
					Id:   uint32(actor.ID),
					Name: actor.Name,
				})
			}
			return protoCasts
		}(),
	}
}

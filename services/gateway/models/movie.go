package models

import (
	"movie/gateway/pb"
)

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
	ReleaseDate uint64
	Director    Director
	Casts       []Actor
}

func DirectorFromProto(p *pb.Director) Director {
	return Director{
		ID:   p.Id,
		Name: p.Name,
	}
}

func ActorFromProto(p *pb.Actor) Actor {
	return Actor{
		ID:   p.Id,
		Name: p.Name,
	}
}

func MovieFromProto(p *pb.Movie) Movie {
	return Movie{
		ID:          p.Id,
		Title:       p.Title,
		Synopsis:    p.Synopsis,
		ReleaseDate: p.ReleaseDate,
		Director: Director{
			ID:   p.Director.Id,
			Name: p.Director.Name,
		},
		Casts: func() []Actor {
			casts := make([]Actor, 0)
			for _, protoActor := range p.Casts {
				casts = append(casts, Actor{
					ID:   protoActor.Id,
					Name: protoActor.Name,
				})
			}
			return casts
		}(),
	}
}

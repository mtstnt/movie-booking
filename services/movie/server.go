package main

import (
	"context"
	"movie/movie/pb"
)

type Server struct {
	pb.UnimplementedMovieServiceServer
	movieRepository Repository
}

func NewServer(movieRepository Repository) Server {
	return Server{
		movieRepository: movieRepository,
	}
}

func (s Server) GetMovies(ctx context.Context, request *pb.GetMoviesRequest) (*pb.GetMoviesResponse, error) {
	movies, err := s.movieRepository.GetMovies()
	if err != nil {
		return nil, err
	}

	protoMovies := make([]*pb.Movie, 0)
	for _, movie := range movies {
		// Experiment with lambdas for encapsulating logic in fields.
		protoMovies = append(protoMovies, MovieToProto(&movie))
	}

	return &pb.GetMoviesResponse{
		Movies: protoMovies,
	}, nil
}

func (s Server) GetMovie(ctx context.Context, request *pb.GetMovieRequest) (*pb.GetMovieResponse, error) {
	id := request.Id

	movie, err := s.movieRepository.GetMovie(id)
	if err != nil {
		return nil, err
	}

	return &pb.GetMovieResponse{
		Movie: &pb.Movie{
			Id:          uint32(movie.ID),
			Title:       movie.Title,
			Synopsis:    movie.Synopsis,
			ReleaseDate: 0,
			Director:    &pb.Director{},
			Casts:       []*pb.Actor{},
		},
	}, nil
}

func (s Server) CreateMovie(ctx context.Context, request *pb.CreateMovieRequest) (*pb.CreateMovieResponse, error) {
	movie := Movie{
		Title:      request.Title,
		Synopsis:   request.Synopsis,
		DirectorID: uint(request.DirectorID),
	}

	if err := s.movieRepository.CreateMovie(&movie); err != nil {
		return nil, err
	}

	return &pb.CreateMovieResponse{
		Movie: MovieToProto(&movie),
	}, nil
}

func (s Server) UpdateMovie(ctx context.Context, request *pb.UpdateMovieRequest) (*pb.UpdateMovieResponse, error) {
	movie, err := s.movieRepository.GetMovie(request.Id)
	if err != nil {
		return nil, err
	}

	if request.Title != nil {
		movie.Title = *request.Title
	}
	if request.Synopsis != nil {
		movie.Synopsis = *request.Synopsis
	}
	if request.ReleaseDate != nil {
		movie.ReleaseDate = uint64(*request.ReleaseDate)
	}

	if err := s.movieRepository.UpdateMovie(&movie); err != nil {
		return nil, err
	}

	return &pb.UpdateMovieResponse{
		Movie: MovieToProto(&movie),
	}, nil
}

func (s Server) DeleteMovie(ctx context.Context, request *pb.DeleteMovieRequest) (*pb.DeleteMovieResponse, error) {
	movie, err := s.movieRepository.GetMovie(request.Id)
	if err != nil {
		return nil, err
	}

	if err := s.movieRepository.DeleteMovie(&movie); err != nil {
		return nil, err
	}

	return &pb.DeleteMovieResponse{}, nil
}

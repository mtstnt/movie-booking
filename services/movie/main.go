package main

import (
	"context"
	"log"
	"movie/gateway/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MovieServer struct {
	pb.UnimplementedMovieServiceServer
}

func (MovieServer) GetMovies(context.Context, pb.GetMoviesRequest) (*pb.GetMoviesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMovies not implemented")
}
func (MovieServer) GetMovie(context.Context, pb.GetMovieRequest) (*pb.GetMovieResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMovie not implemented")
}
func (MovieServer) CreateMovie(context.Context, pb.CreateMovieRequest) (*pb.CreateMovieResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateMovie not implemented")
}
func (MovieServer) UpdateMovie(context.Context, pb.UpdateMovieRequest) (*pb.UpdateMovieResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMovie not implemented")
}
func (MovieServer) DeleteMovie(context.Context, pb.DeleteMovieRequest) (*pb.DeleteMovieResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteMovie not implemented")
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	// listener, err := net.Listen("tcp", "0.0.0.0:5000")
	// if err != nil {
	// 	return err
	// }

	// grpcServer, err := "", errors.New("")

	return nil
}

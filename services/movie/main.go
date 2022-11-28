package main

import (
	"log"
	"movie/movie/pb"
	"net"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	listener, err := net.Listen("tcp", "0.0.0.0:5000")
	if err != nil {
		return err
	}

	db, err := gorm.Open(
		postgres.Open("host=movies_db user=movies password=movies dbname=movies port=5432 TimeZone=Asia/Jakarta"),
	)
	if err != nil {
		return err
	}

	if err := db.AutoMigrate(Director{}, Actor{}, Movie{}); err != nil {
		return err
	}

	movieRepository := NewRepository(db)
	movieServer := NewServer(movieRepository)

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	pb.RegisterMovieServiceServer(grpcServer, movieServer)

	log.Println("Running GRPC Server at port 8000")
	return grpcServer.Serve(listener)
}

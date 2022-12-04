package main

import (
	"log"
	"movie/user/pb"
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

func resetDB(db *gorm.DB) error {
	return db.Migrator().DropTable("users")
}

func run() error {
	listener, err := net.Listen("tcp", "0.0.0.0:8000")
	if err != nil {
		return err
	}

	db, err := gorm.Open(
		postgres.Open("host=users_db user=users password=users dbname=users port=5432"),
	)
	if err != nil {
		return err
	}

	if err := resetDB(db); err != nil {
		return err
	}
	if err := db.AutoMigrate(User{}); err != nil {
		return err
	}

	userRepository := NewRepository(db)
	userService := NewService(userRepository)
	userServer := NewServer(userService, userRepository)

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	pb.RegisterUserServiceServer(grpcServer, userServer)

	log.Println("Running GRPC Server at port 8000")
	return grpcServer.Serve(listener)
}

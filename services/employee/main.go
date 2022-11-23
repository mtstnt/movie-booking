package main

import (
	"log"
	"movie/employee/pb"
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
	listener, err := net.Listen("tcp", "0.0.0.0:8000")
	if err != nil {
		return err
	}

	db, err := gorm.Open(
		postgres.Open("host=employees_db user=employees password=employees dbname=employees port=5432 TimeZone=Asia/Jakarta"),
	)
	if err != nil {
		return err
	}

	if err := db.AutoMigrate(Employee{}); err != nil {
		return err
	}

	employeeRepository := NewRepository(db)
	employeeService := NewService(employeeRepository)
	employeeServer := NewServer(employeeService, employeeRepository)

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	pb.RegisterEmployeeServiceServer(grpcServer, employeeServer)

	log.Println("Running GRPC Server at port 8000")
	return grpcServer.Serve(listener)
}

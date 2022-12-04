package main

import (
	"log"
	"movie/booking/model"
	"movie/booking/pb"
	"movie/booking/repository"
	"movie/booking/server"
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
	return db.Migrator().DropTable("schedules", "bookings")
}

func run() error {
	listener, err := net.Listen("tcp", "0.0.0.0:8000")
	if err != nil {
		return err
	}

	db, err := gorm.Open(
		postgres.Open("host=bookings_db user=bookings password=bookings dbname=bookings port=5432"),
	)
	if err != nil {
		return err
	}

	if err := resetDB(db); err != nil {
		return err
	}
	if err := db.AutoMigrate(model.Schedule{}, model.Booking{}); err != nil {
		return err
	}

	scheduleRepository := repository.NewScheduleRepo(db)
	scheduleServer := server.NewScheduleServer(scheduleRepository)

	bookingRepository := repository.NewBookingRepo(db)
	bookingServer := server.NewBookingServer(bookingRepository, scheduleRepository)

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	pb.RegisterBookingServiceServer(grpcServer, bookingServer)
	pb.RegisterScheduleServiceServer(grpcServer, scheduleServer)

	log.Println("Running GRPC Server at port 8000")
	return grpcServer.Serve(listener)
}

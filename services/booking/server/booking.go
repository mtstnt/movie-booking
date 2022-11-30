package server

import (
	"context"
	"movie/booking/model"
	"movie/booking/pb"
	"movie/booking/repository"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BookingServer struct {
	pb.UnimplementedBookingServiceServer

	bookingRepo repository.BookingRepo
}

func NewBookingServer(bookingRepo repository.BookingRepo) BookingServer {
	return BookingServer{
		bookingRepo: bookingRepo,
	}
}

func (s BookingServer) GetUserBookings(ctx context.Context, request *pb.GetUserBookingsRequest) (*pb.GetUserBookingsResponse, error) {
	bookings, err := s.bookingRepo.GetUserBookings(request.UserID, request.From, request.To, request.MovieID)
	if err != nil {
		return nil, err
	}

	protoBookings := make([]*pb.Booking, 0)
	for _, booking := range bookings {
		protoBookings = append(protoBookings, model.BookingToProto(&booking))
	}

	return &pb.GetUserBookingsResponse{
		Bookings: protoBookings,
	}, nil
}

func (s BookingServer) CreateBooking(context.Context, *pb.CreateBookingRequest) (*pb.CreateBookingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateBooking not implemented")
}

func (s BookingServer) CancelBooking(context.Context, *pb.CancelBookingRequest) (*pb.CancelBookingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelBooking not implemented")
}

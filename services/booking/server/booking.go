package server

import (
	"context"
	"movie/booking/model"
	"movie/booking/pb"
	"movie/booking/repository"
)

type BookingServer struct {
	pb.UnimplementedBookingServiceServer

	bookingRepo  repository.BookingRepo
	scheduleRepo repository.ScheduleRepo
}

func NewBookingServer(
	bookingRepo repository.BookingRepo,
	scheduleRepo repository.ScheduleRepo,
) BookingServer {
	return BookingServer{
		bookingRepo:  bookingRepo,
		scheduleRepo: scheduleRepo,
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

func (s BookingServer) CreateBooking(ctx context.Context, request *pb.CreateBookingRequest) (*pb.CreateBookingResponse, error) {
	schedule, err := s.scheduleRepo.GetSchedule(request.ScheduleID)
	if err != nil {
		return nil, err
	}

	booking := model.Booking{
		UserID:     request.UserID,
		ScheduleID: uint32(schedule.ID),
		IsCanceled: false,
	}

	if err := s.bookingRepo.CreateBooking(&booking); err != nil {
		return nil, err
	}

	return &pb.CreateBookingResponse{
		Booking: model.BookingToProto(&booking),
	}, nil
}

func (s BookingServer) CancelBooking(ctx context.Context, request *pb.CancelBookingRequest) (*pb.CancelBookingResponse, error) {
	booking, err := s.bookingRepo.GetBooking(request.UserID, request.BookingID)
	if err != nil {
		return nil, err
	}

	if err := s.bookingRepo.DeleteBooking(&booking); err != nil {
		return nil, err
	}

	return &pb.CancelBookingResponse{}, nil
}

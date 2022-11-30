package server

import (
	"context"
	"movie/booking/model"
	"movie/booking/pb"
	"movie/booking/repository"
)

type ScheduleServer struct {
	pb.UnimplementedScheduleServiceServer
	scheduleRepo repository.ScheduleRepo
}

func NewScheduleServer(scheduleRepo repository.ScheduleRepo) ScheduleServer {
	return ScheduleServer{
		scheduleRepo: scheduleRepo,
	}
}

func (s ScheduleServer) GetSchedules(ctx context.Context, request *pb.GetSchedulesRequest) (*pb.GetSchedulesResponse, error) {
	schedules, err := s.scheduleRepo.GetAllSchedules(request.From, request.To, request.MovieID)
	if err != nil {
		return nil, err
	}

	protoSchedules := make([]*pb.Schedule, 0)
	for _, schedule := range schedules {
		protoSchedules = append(protoSchedules, model.ScheduleToProto(&schedule))
	}

	return &pb.GetSchedulesResponse{
		Schedules: protoSchedules,
	}, nil
}

func (s ScheduleServer) GetSchedule(ctx context.Context, request *pb.GetScheduleRequest) (*pb.GetScheduleResponse, error) {
	id := request.Id

	schedule, err := s.scheduleRepo.GetSchedule(id)
	if err != nil {
		return nil, err
	}

	return &pb.GetScheduleResponse{
		Schedule: model.ScheduleToProto(&schedule),
	}, nil
}

func (s ScheduleServer) CreateSchedule(ctx context.Context, request *pb.CreateScheduleRequest) (*pb.CreateScheduleResponse, error) {
	schedule := model.Schedule{
		MovieID:  request.MovieID,
		ShowTime: request.ShowTime,
		StudioNo: request.StudioNo,
		Capacity: request.Capacity,
	}

	if err := s.scheduleRepo.CreateSchedule(&schedule); err != nil {
		return nil, err
	}

	return &pb.CreateScheduleResponse{
		Schedule: model.ScheduleToProto(&schedule),
	}, nil
}

func (s ScheduleServer) UpdateSchedule(ctx context.Context, request *pb.UpdateScheduleRequest) (*pb.UpdateScheduleResponse, error) {
	schedule, err := s.scheduleRepo.GetSchedule(request.Id)
	if err != nil {
		return nil, err
	}

	if request.MovieID != nil {
		schedule.MovieID = *request.MovieID
	}

	if err := s.scheduleRepo.UpdateSchedule(&schedule); err != nil {
		return nil, err
	}

	return &pb.UpdateScheduleResponse{
		Schedule: model.ScheduleToProto(&schedule),
	}, nil
}

func (s ScheduleServer) DeleteSchedule(ctx context.Context, request *pb.DeleteScheduleRequest) (*pb.DeleteScheduleResponse, error) {
	schedule, err := s.scheduleRepo.GetSchedule(request.Id)
	if err != nil {
		return nil, err
	}

	if err := s.scheduleRepo.DeleteSchedule(&schedule); err != nil {
		return nil, err
	}

	return &pb.DeleteScheduleResponse{}, nil
}

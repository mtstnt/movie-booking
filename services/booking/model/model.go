package model

import (
	"movie/booking/pb"

	"gorm.io/gorm"
)

type Schedule struct {
	gorm.Model
	ShowTime uint64
	MovieID  uint32
	StudioNo uint32
	Capacity uint32
}

type Booking struct {
	gorm.Model
	UserID     uint32
	ScheduleID uint32
	IsCanceled bool
}

func ScheduleFromProto(s *pb.Schedule) Schedule {
	return Schedule{
		MovieID:  s.MovieID,
		ShowTime: s.ShowTime,
		StudioNo: s.StudioNo,
		Capacity: s.Capacity,
	}
}

func ScheduleToProto(s *Schedule) *pb.Schedule {
	return &pb.Schedule{
		Id:       uint32(s.ID),
		MovieID:  s.MovieID,
		ShowTime: s.ShowTime,
		StudioNo: s.StudioNo,
		Capacity: s.Capacity,
	}
}

func BookingFromProto(b *pb.Booking) Booking {
	return Booking{
		Model: gorm.Model{
			ID: uint(b.Id),
		},
		UserID:     b.UserID,
		ScheduleID: b.ScheduleID,
		IsCanceled: b.IsCanceled,
	}
}

func BookingToProto(b *Booking) *pb.Booking {
	return &pb.Booking{
		Id:         uint32(b.ID),
		UserID:     b.UserID,
		ScheduleID: b.ScheduleID,
		IsCanceled: b.IsCanceled,
	}
}

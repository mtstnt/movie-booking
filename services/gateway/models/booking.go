package models

import "movie/gateway/pb"

type Booking struct {
	ID         uint32
	UserID     uint32
	ScheduleID uint32
	IsCanceled bool
}

type Schedule struct {
	ID       uint32
	MovieID  uint32
	ShowTime uint64
	StudioNo uint32
	Capacity uint32
}

func BookingFromProto(b *pb.Booking) Booking {
	return Booking{
		ID:         b.Id,
		UserID:     b.UserID,
		ScheduleID: b.ScheduleID,
		IsCanceled: b.IsCanceled,
	}
}

func ScheduleFromProto(b *pb.Schedule) Schedule {
	return Schedule{
		ID:       b.Id,
		MovieID:  b.MovieID,
		ShowTime: b.ShowTime,
		StudioNo: b.StudioNo,
		Capacity: b.Capacity,
	}
}

package repository

import (
	"movie/booking/model"
	"time"

	"gorm.io/gorm"
)

type ScheduleRepo struct {
	db *gorm.DB
}

func NewScheduleRepo(db *gorm.DB) ScheduleRepo {
	return ScheduleRepo{db}
}

func (r ScheduleRepo) GetAllSchedules(from uint64, to uint64, movieID uint32) ([]model.Schedule, error) {
	var schedule []model.Schedule
	if to == 0 {
		to = uint64(time.Now().Day())
	}

	tx := r.db
	if movieID > 0 {
		tx = tx.Where("movie_id = ?", movieID)
	}
	if err := tx.Where("show_time > ? AND show_time < ?", from, to).Find(&schedule).Error; err != nil {
		return nil, err
	}
	return schedule, nil
}

func (r ScheduleRepo) GetSchedule(id uint32) (model.Schedule, error) {
	var schedule model.Schedule
	if err := r.db.First(&schedule, "id = ?", id).Error; err != nil {
		return schedule, err
	}
	return schedule, nil
}

func (r ScheduleRepo) CreateSchedule(schedule *model.Schedule) error {
	return r.db.Create(schedule).Error
}

func (r ScheduleRepo) UpdateSchedule(schedule *model.Schedule) error {
	return r.db.Updates(schedule).Error
}

func (r ScheduleRepo) DeleteSchedule(schedule *model.Schedule) error {
	schedule.DeletedAt = gorm.DeletedAt{
		Time:  time.Now(),
		Valid: true,
	}
	return r.db.Updates(schedule).Error
}

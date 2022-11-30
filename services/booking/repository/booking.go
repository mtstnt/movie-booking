package repository

import (
	"movie/booking/model"
	"time"

	"gorm.io/gorm"
)

type BookingRepo struct {
	db *gorm.DB
}

func NewBookingRepo(db *gorm.DB) BookingRepo {
	return BookingRepo{db}
}

func (r BookingRepo) GetUserBookings(userID uint32, from uint64, to uint64, movieID uint32) ([]model.Booking, error) {
	var bookings []model.Booking

	tx := r.db.Where("user_id = ?", userID)
	tx = tx.Where("show_time > ?", from)

	if to == 0 {
		now := time.Now()
		location, err := time.LoadLocation("Asia/Jakarta")

		if err != nil {
			return nil, err
		}

		to = uint64(time.Date(
			now.Year(),
			now.Month(),
			now.Day(),
			23, 59, 0, 0,
			location,
		).Unix())

		tx = tx.Where("show_time < ?", to)
	}

	if movieID > 0 {
		tx = tx.Where("movie_id = ?", movieID)
	}

	if err := tx.Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r BookingRepo) GetBooking(userId uint32, id uint32) (model.Booking, error) {
	var booking model.Booking
	if err := r.db.First(&booking, "id = ?", id).Error; err != nil {
		return booking, err
	}
	return booking, nil
}

func (r BookingRepo) CreateBooking(booking *model.Booking) error {
	return r.db.Create(booking).Error
}

func (r BookingRepo) UpdateBooking(booking *model.Booking) error {
	return r.db.Updates(booking).Error
}

func (r BookingRepo) DeleteBooking(booking *model.Booking) error {
	booking.DeletedAt = gorm.DeletedAt{
		Time:  time.Now(),
		Valid: true,
	}
	return r.db.Updates(booking).Error
}

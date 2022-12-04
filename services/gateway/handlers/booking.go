package handlers

import (
	"movie/gateway/helpers"
	"movie/gateway/models"
	"movie/gateway/pb"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nleeper/goment"
)

type Booking struct {
	bookingClient  pb.BookingServiceClient
	scheduleClient pb.ScheduleServiceClient
}

func NewBooking(
	bookingClient pb.BookingServiceClient,
	scheduleClient pb.ScheduleServiceClient,
) Booking {
	return Booking{bookingClient, scheduleClient}
}

func (b Booking) GetUserBookings(ctx *gin.Context) {
	user := ctx.Value("User").(models.User)

	formattedFrom, formattedTo := ctx.Query("from"), ctx.Query("to")
	movie := ctx.Query("movie")

	var unixFrom, unixTo int64 = 0, 0
	var movieID uint32 = 0

	if formattedFrom != "" {
		fromTimeObj, err := goment.New(formattedFrom)
		if err != nil {
			helpers.HttpError(ctx, http.StatusBadRequest, err)
			return
		}
		unixFrom = fromTimeObj.ToUnix()
	}

	if formattedTo != "" {
		toTimeObj, err := goment.New(formattedTo)
		if err != nil {
			helpers.HttpError(ctx, http.StatusBadRequest, err)
			return
		}
		unixTo = toTimeObj.ToUnix()
	}

	if movie != "" {
		movieIDTmp, err := strconv.ParseUint(movie, 10, 32)
		if err != nil {
			helpers.HttpError(ctx, http.StatusBadRequest, err)
			return
		}
		movieID = uint32(movieIDTmp)
	}

	response, err := b.bookingClient.GetUserBookings(ctx, &pb.GetUserBookingsRequest{
		UserID:  user.ID,
		From:    uint64(unixFrom),
		To:      uint64(unixTo),
		MovieID: movieID,
	})
	if err != nil {
		helpers.HttpError(ctx, http.StatusInternalServerError, err)
		return
	}

	bookings := make([]models.Booking, 0)
	for _, protoBooking := range response.Bookings {
		bookings = append(bookings, models.BookingFromProto(protoBooking))
	}

	helpers.HttpOK(ctx, gin.H{
		"Bookings": bookings,
	})
}

func (b Booking) GetBooking(ctx *gin.Context) {
	user := ctx.Value("User").(models.User)

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		helpers.HttpError(ctx, http.StatusBadRequest, err)
		return
	}

	response, err := b.bookingClient.GetBooking(ctx, &pb.GetBookingRequest{
		Id:     uint32(id),
		UserID: user.ID,
	})
	if err != nil {
		helpers.HttpError(ctx, http.StatusInternalServerError, err)
		return
	}

	booking := models.BookingFromProto(response.Booking)

	helpers.HttpOK(ctx, gin.H{
		"Booking": booking,
	})
}

func (b Booking) CreateBooking(ctx *gin.Context) {
	user := ctx.Value("User").(models.User)

	var req struct {
		ScheduleID uint32
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.HttpError(ctx, http.StatusBadRequest, err)
		return
	}

	response, err := b.bookingClient.CreateBooking(ctx, &pb.CreateBookingRequest{
		UserID:     user.ID,
		ScheduleID: req.ScheduleID,
	})
	if err != nil {
		helpers.HttpError(ctx, http.StatusInternalServerError, err)
		return
	}

	booking := models.BookingFromProto(response.Booking)

	helpers.HttpOK(ctx, gin.H{
		"Booking": booking,
	})
}

func (b Booking) CancelBooking(ctx *gin.Context) {
	user := ctx.Value("User").(models.User)

	var req struct {
		BookingID uint32
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.HttpError(ctx, http.StatusBadRequest, err)
		return
	}

	_, err := b.bookingClient.CancelBooking(ctx, &pb.CancelBookingRequest{
		BookingID: req.BookingID,
		UserID:    user.ID,
	})
	if err != nil {
		helpers.HttpError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.AbortWithStatus(http.StatusNoContent)
}

func (b Booking) GetSchedules(ctx *gin.Context) {
	var from, to uint64 = 0, 0
	var movieID uint64 = 0
	var fromStr, toStr string

	fromStr = ctx.Query("from")
	if fromStr != "" {
		var err error
		from, err = strconv.ParseUint(fromStr, 10, 64)
		if err != nil {
			helpers.HttpError(ctx, http.StatusBadRequest, err)
			return
		}
	}

	toStr = ctx.Query("to")
	if toStr != "" {
		var err error
		to, err = strconv.ParseUint(toStr, 10, 64)
		if err != nil {
			helpers.HttpError(ctx, http.StatusBadRequest, err)
			return
		}
	}

	movieIDStr := ctx.Query("movie")
	if movieIDStr != "" {
		var err error
		movieID, err = strconv.ParseUint(movieIDStr, 10, 32)
		if err != nil {
			helpers.HttpError(ctx, http.StatusBadRequest, err)
			return
		}
	}

	response, err := b.scheduleClient.GetSchedules(ctx, &pb.GetSchedulesRequest{
		From:    from,
		To:      to,
		MovieID: uint32(movieID),
	})
	if err != nil {
		helpers.HttpError(ctx, http.StatusInternalServerError, err)
		return
	}

	schedules := []models.Schedule{}
	for _, protoSchedule := range response.Schedules {
		schedules = append(schedules, models.ScheduleFromProto(protoSchedule))
	}

	helpers.HttpOK(ctx, gin.H{
		"Schedules": schedules,
	})
}

func (b Booking) GetSchedule(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		helpers.HttpError(ctx, http.StatusBadRequest, err)
		return
	}

	response, err := b.scheduleClient.GetSchedule(ctx, &pb.GetScheduleRequest{
		Id: uint32(id),
	})
	if err != nil {
		helpers.HttpError(ctx, http.StatusInternalServerError, err)
		return
	}

	helpers.HttpOK(ctx, gin.H{
		"Schedule": models.ScheduleFromProto(response.Schedule),
	})
}

func (b Booking) CreateSchedule(ctx *gin.Context) {
	var req struct {
		MovieID  uint32
		ShowTime uint64
		Capacity uint32
		StudioNo uint32
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.HttpError(ctx, http.StatusBadRequest, err)
		return
	}

	response, err := b.scheduleClient.CreateSchedule(ctx, &pb.CreateScheduleRequest{
		MovieID:  req.MovieID,
		StudioNo: req.StudioNo,
		Capacity: req.Capacity,
		ShowTime: req.ShowTime,
	})
	if err != nil {
		helpers.HttpError(ctx, http.StatusInternalServerError, err)
		return
	}

	helpers.HttpOK(ctx, gin.H{
		"Schedule": models.ScheduleFromProto(response.Schedule),
	})
}

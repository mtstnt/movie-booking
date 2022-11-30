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
	bookingClient pb.BookingServiceClient
}

func NewBooking(bookingClient pb.BookingServiceClient) Booking {
	return Booking{bookingClient}
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

}

func (b Booking) CreateBooking(ctx *gin.Context) {

}

func (b Booking) UpdateBooking(ctx *gin.Context) {

}

func (b Booking) CancelBooking(ctx *gin.Context) {

}

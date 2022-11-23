package handlers

import (
	"movie/gateway/pb"

	"github.com/gin-gonic/gin"
)

type Movie struct {
	employeeService pb.EmployeeServiceClient
	movieService    pb.MovieServiceClient
}

func NewMovie(
	employeeService pb.EmployeeServiceClient,
	movieService pb.MovieServiceClient,
) Movie {
	return Movie{employeeService, movieService}
}

func (m Movie) GetMovies(ctx *gin.Context) {

}

func (m Movie) GetMovie(ctx *gin.Context) {

}

func (m Movie) CreateMovie(ctx *gin.Context) {

}

func (m Movie) UpdateMovie(ctx *gin.Context) {

}

func (m Movie) DeleteMovie(ctx *gin.Context) {

}

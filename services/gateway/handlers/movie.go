package handlers

import (
	"movie/gateway/helpers"
	"movie/gateway/models"
	"movie/gateway/pb"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Movie struct {
	employeeClient pb.EmployeeServiceClient
	movieClient    pb.MovieServiceClient
}

func NewMovie(
	employeeClient pb.EmployeeServiceClient,
	movieClient pb.MovieServiceClient,
) Movie {
	return Movie{employeeClient, movieClient}
}

func (m Movie) GetMovies(ctx *gin.Context) {
	response, err := m.movieClient.GetMovies(ctx, &pb.GetMoviesRequest{})
	if err != nil {
		helpers.HttpError(ctx, http.StatusInternalServerError, err)
		return
	}

	movies := make([]models.Movie, 0)
	for _, protoMovie := range response.Movies {
		movies = append(movies, models.MovieFromProto(protoMovie))
	}

	helpers.HttpOK(ctx, gin.H{
		"Movies": movies,
	})
}

func (m Movie) GetMovie(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		helpers.HttpError(ctx, http.StatusBadRequest, err)
		return
	}

	response, err := m.movieClient.GetMovie(ctx, &pb.GetMovieRequest{
		Id: uint32(id),
	})
	if err != nil {
		helpers.HttpError(ctx, http.StatusBadRequest, err)
	}

	helpers.HttpOK(ctx, gin.H{
		"Movie": models.MovieFromProto(response.Movie),
	})
}

func (m Movie) CreateMovie(ctx *gin.Context) {
	var req struct {
		Title       string
		Synopsis    string
		ReleaseDate uint64
		DirectorID  uint32
		CastsID     []uint32
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.HttpError(ctx, http.StatusBadRequest, err)
		return
	}

	response, err := m.movieClient.CreateMovie(ctx, &pb.CreateMovieRequest{
		Title:       req.Title,
		Synopsis:    req.Synopsis,
		ReleaseDate: req.ReleaseDate,
		DirectorID:  req.DirectorID,
		CastsID:     req.CastsID,
	})
	if err != nil {
		helpers.HttpError(ctx, http.StatusInternalServerError, err)
		return
	}

	helpers.HttpOK(ctx, gin.H{
		"Movie": models.MovieFromProto(response.Movie),
	})
}

func (m Movie) UpdateMovie(ctx *gin.Context) {
	var req struct {
		Title       string   `json:"Title,omitempty"`
		Synopsis    string   `json:"Synopsis,omitempty"`
		ReleaseDate uint64   `json:"ReleaseDate,omitempty"`
		DirectorID  uint32   `json:"DirectorID,omitempty"`
		CastsID     []uint32 `json:"CastsID,omitempty"`
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		helpers.HttpError(ctx, http.StatusBadRequest, err)
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.HttpError(ctx, http.StatusBadRequest, err)
		return
	}

	grpcReq := &pb.UpdateMovieRequest{Id: uint32(id)}
	if req.Title != "" {
		grpcReq.Title = &req.Title
	}
	if req.Synopsis != "" {
		grpcReq.Synopsis = &req.Synopsis
	}
	if req.ReleaseDate != 0 {
		grpcReq.ReleaseDate = &req.ReleaseDate
	}
	if req.DirectorID != 0 {
		grpcReq.DirectorID = &req.DirectorID
	}
	if len(req.CastsID) != 0 {
		grpcReq.CastsID = req.CastsID
	}

	response, err := m.movieClient.UpdateMovie(ctx, grpcReq)
	if err != nil {
		helpers.HttpError(ctx, http.StatusInternalServerError, err)
		return
	}

	helpers.HttpOK(ctx, gin.H{
		"Movie": models.MovieFromProto(response.Movie),
	})
}

func (m Movie) DeleteMovie(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		helpers.HttpError(ctx, http.StatusBadRequest, err)
		return
	}

	if _, err = m.movieClient.DeleteMovie(ctx, &pb.DeleteMovieRequest{Id: uint32(id)}); err != nil {
		helpers.HttpError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.AbortWithStatus(http.StatusNoContent)
}

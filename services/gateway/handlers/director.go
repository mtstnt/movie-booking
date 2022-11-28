package handlers

import (
	"movie/gateway/helpers"
	"movie/gateway/models"
	"movie/gateway/pb"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Director struct {
	directorClient pb.DirectorServiceClient
}

func (m Director) GetDirectors(ctx *gin.Context) {
	response, err := m.directorClient.GetDirectors(ctx, &pb.GetDirectorsRequest{})
	if err != nil {
		helpers.HttpError(ctx, http.StatusInternalServerError, err)
		return
	}

	directors := make([]models.Director, 0)
	for _, protoDirector := range response.Directors {
		directors = append(directors, models.DirectorFromProto(protoDirector))
	}

	helpers.HttpOK(ctx, gin.H{
		"Directors": directors,
	})
}

func (m Director) GetDirector(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		helpers.HttpError(ctx, http.StatusBadRequest, err)
		return
	}

	response, err := m.directorClient.GetDirector(ctx, &pb.GetDirectorRequest{
		Id: uint32(id),
	})
	if err != nil {
		helpers.HttpError(ctx, http.StatusBadRequest, err)
	}

	helpers.HttpOK(ctx, gin.H{
		"Director": models.DirectorFromProto(response.Director),
	})
}

func (m Director) CreateDirector(ctx *gin.Context) {
	var req struct {
		Name string
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.HttpError(ctx, http.StatusBadRequest, err)
		return
	}

	response, err := m.directorClient.CreateDirector(ctx, &pb.CreateDirectorRequest{
		Name: req.Name,
	})
	if err != nil {
		helpers.HttpError(ctx, http.StatusInternalServerError, err)
		return
	}

	helpers.HttpOK(ctx, gin.H{
		"Director": models.DirectorFromProto(response.Director),
	})
}

func (m Director) UpdateDirector(ctx *gin.Context) {
	var req struct {
		Name string `json:"Name,omitempty"`
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

	grpcReq := &pb.UpdateDirectorRequest{Id: uint32(id)}
	if req.Name != "" {
		grpcReq.Name = &req.Name
	}

	response, err := m.directorClient.UpdateDirector(ctx, grpcReq)
	if err != nil {
		helpers.HttpError(ctx, http.StatusInternalServerError, err)
		return
	}

	helpers.HttpOK(ctx, gin.H{
		"Director": models.DirectorFromProto(response.Director),
	})
}

func (m Director) DeleteDirector(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		helpers.HttpError(ctx, http.StatusBadRequest, err)
		return
	}

	if _, err = m.directorClient.DeleteDirector(ctx, &pb.DeleteDirectorRequest{Id: uint32(id)}); err != nil {
		helpers.HttpError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.AbortWithStatus(http.StatusNoContent)
}

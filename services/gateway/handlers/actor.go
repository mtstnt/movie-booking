package handlers

import (
	"movie/gateway/helpers"
	"movie/gateway/models"
	"movie/gateway/pb"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Actor struct {
	actorClient pb.ActorServiceClient
}

func (m Actor) GetActors(ctx *gin.Context) {
	response, err := m.actorClient.GetActors(ctx, &pb.GetActorsRequest{})
	if err != nil {
		helpers.HttpError(ctx, http.StatusInternalServerError, err)
		return
	}

	actors := make([]models.Actor, 0)
	for _, protoActor := range response.Actors {
		actors = append(actors, models.ActorFromProto(protoActor))
	}

	helpers.HttpOK(ctx, gin.H{
		"Actors": actors,
	})
}

func (m Actor) GetActor(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		helpers.HttpError(ctx, http.StatusBadRequest, err)
		return
	}

	response, err := m.actorClient.GetActor(ctx, &pb.GetActorRequest{
		Id: uint32(id),
	})
	if err != nil {
		helpers.HttpError(ctx, http.StatusBadRequest, err)
	}

	helpers.HttpOK(ctx, gin.H{
		"Actor": models.ActorFromProto(response.Actor),
	})
}

func (m Actor) CreateActor(ctx *gin.Context) {
	var req struct {
		Name string
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.HttpError(ctx, http.StatusBadRequest, err)
		return
	}

	response, err := m.actorClient.CreateActor(ctx, &pb.CreateActorRequest{
		Name: req.Name,
	})
	if err != nil {
		helpers.HttpError(ctx, http.StatusInternalServerError, err)
		return
	}

	helpers.HttpOK(ctx, gin.H{
		"Actor": models.ActorFromProto(response.Actor),
	})
}

func (m Actor) UpdateActor(ctx *gin.Context) {
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

	grpcReq := &pb.UpdateActorRequest{Id: uint32(id)}
	if req.Name != "" {
		grpcReq.Name = &req.Name
	}

	response, err := m.actorClient.UpdateActor(ctx, grpcReq)
	if err != nil {
		helpers.HttpError(ctx, http.StatusInternalServerError, err)
		return
	}

	helpers.HttpOK(ctx, gin.H{
		"Actor": models.ActorFromProto(response.Actor),
	})
}

func (m Actor) DeleteActor(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		helpers.HttpError(ctx, http.StatusBadRequest, err)
		return
	}

	if _, err = m.actorClient.DeleteActor(ctx, &pb.DeleteActorRequest{Id: uint32(id)}); err != nil {
		helpers.HttpError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.AbortWithStatus(http.StatusNoContent)
}

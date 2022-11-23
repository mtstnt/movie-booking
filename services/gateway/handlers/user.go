package handlers

import (
	"movie/gateway/helpers"
	"movie/gateway/models"
	"movie/gateway/pb"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
)

type User struct {
	userClient  pb.UserServiceClient
	redisClient *redis.Client
}

func NewUser(
	userClient pb.UserServiceClient,
	redisClient *redis.Client,
) User {
	return User{userClient, redisClient}
}

func (u User) SignIn(ctx *gin.Context) {
	var req struct {
		Email    string
		Password string
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.HttpError(ctx, http.StatusBadRequest, err)
		return
	}

	response, err := u.userClient.AuthenticateUser(ctx.Request.Context(), &pb.AuthenticateUserRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		helpers.HttpError(ctx, http.StatusUnauthorized, err)
		return
	}

	user := models.User{
		ID:    response.User.Id,
		Email: response.User.Email,
		Name:  response.User.Name,
	}

	token, err := helpers.CreateUserToken(&user)
	if err != nil {
		helpers.HttpError(ctx, http.StatusInternalServerError, err)
		return
	}

	u.redisClient.Set(ctx.Request.Context(), token, user.ID, 1*time.Hour)

	helpers.HttpOK(ctx, gin.H{
		"User":  user,
		"Token": token,
	})
}

func (u User) SignUp(ctx *gin.Context) {
	var req struct {
		Email    string
		Name     string
		Password string
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.HttpError(ctx, http.StatusBadRequest, err)
		return
	}

	response, err := u.userClient.CreateUser(ctx.Request.Context(), &pb.CreateUserRequest{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
	})
	if err != nil {
		helpers.HttpError(ctx, http.StatusInternalServerError, err)
		return
	}

	user := models.User{
		ID:    response.User.Id,
		Email: response.User.Email,
		Name:  response.User.Name,
	}

	token, err := helpers.CreateUserToken(&user)
	if err != nil {
		helpers.HttpError(ctx, http.StatusInternalServerError, err)
		return
	}

	u.redisClient.Set(ctx.Request.Context(), token, user.ID, 1*time.Hour)

	helpers.HttpOK(ctx, gin.H{
		"User":  user,
		"Token": token,
	})
}

func (u User) GetAllUsers(ctx *gin.Context) {
	response, err := u.userClient.GetUsers(ctx, &pb.GetUsersRequest{})
	if err != nil {
		helpers.HttpError(ctx, http.StatusInternalServerError, err)
		return
	}

	users := []models.User{}
	for _, protoUser := range response.Users {
		users = append(users, models.User{
			ID:    protoUser.Id,
			Email: protoUser.Email,
			Name:  protoUser.Name,
		})
	}

	helpers.HttpOK(ctx, gin.H{
		"Users": users,
	})
}

func (u User) GetUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		helpers.HttpError(ctx, http.StatusBadRequest, err)
		return
	}

	response, err := u.userClient.GetUser(ctx, &pb.GetUserRequest{
		Id: uint32(id),
	})
	if err != nil {
		helpers.HttpError(ctx, http.StatusBadRequest, err)
	}

	// TODO: Get history booking for user.

	helpers.HttpOK(ctx, gin.H{
		"User": models.User{
			ID:    response.User.Id,
			Email: response.User.Email,
			Name:  response.User.Name,
		},
	})
}

func (u User) UpdateUser(ctx *gin.Context) {
	var req struct {
		Email string `json:"Email,omitempty"`
		Name  string `json:"Name,omitempty"`
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

	grpcRequest := &pb.UpdateUserRequest{}
	grpcRequest.Id = uint32(id)
	if req.Email != "" {
		grpcRequest.Email = &req.Email
	}
	if req.Name != "" {
		grpcRequest.Name = &req.Name
	}

	response, err := u.userClient.UpdateUser(ctx, grpcRequest)
	if err != nil {
		helpers.HttpError(ctx, http.StatusInternalServerError, err)
		return
	}

	helpers.HttpOK(ctx, gin.H{
		"User": models.User{
			ID:    response.User.Id,
			Email: response.User.Email,
			Name:  response.User.Name,
		},
	})
}

func (u User) DeleteUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		helpers.HttpError(ctx, http.StatusBadRequest, err)
		return
	}

	if _, err = u.userClient.DeleteUser(ctx, &pb.DeleteUserRequest{Id: uint32(id)}); err != nil {
		helpers.HttpError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.AbortWithStatus(http.StatusNoContent)
}

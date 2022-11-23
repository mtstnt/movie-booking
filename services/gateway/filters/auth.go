package filters

import (
	"errors"
	"movie/gateway/helpers"
	"movie/gateway/models"
	"movie/gateway/pb"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/golang-jwt/jwt/v4"
)

func getAuthorizationToken(ctx *gin.Context) (string, error) {
	var token string
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return token, errors.New("empty authorization header")
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) < 2 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return token, errors.New("invalid authorization header")
	}

	token = authHeaderParts[0]
	return token, nil
}

func getIDFromJWT(token string) (uint32, error) {
	var claims helpers.Claims
	_, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(helpers.SECRET), nil
	})
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

func AuthenticatedAsUser(userService pb.UserServiceClient, redisClient *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := getAuthorizationToken(ctx)
		if err != nil {
			helpers.HttpError(ctx, http.StatusBadRequest, err)
			return
		}

		redisCmd := redisClient.Get(ctx, token)
		if redisCmd.Err() != nil {
			helpers.HttpError(ctx, http.StatusUnauthorized, err)
			return
		}

		userID, err := getIDFromJWT(token)
		if err != nil {
			helpers.HttpError(ctx, http.StatusUnauthorized, err)
			return
		}

		response, err := userService.GetUser(ctx.Request.Context(), &pb.GetUserRequest{
			Id: userID,
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

		ctx.Set("User", user)
		ctx.Next()
	}
}

func AuthenticatedAsEmployee(employeeService pb.EmployeeServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := getAuthorizationToken(ctx)
		if err != nil {
			helpers.HttpError(ctx, http.StatusBadRequest, err)
			return
		}

		userID, err := getIDFromJWT(token)
		if err != nil {
			helpers.HttpError(ctx, http.StatusUnauthorized, err)
			return
		}

		response, err := employeeService.GetEmployee(ctx.Request.Context(), &pb.GetEmployeeRequest{
			Id: userID,
		})
		if err != nil {
			helpers.HttpError(ctx, http.StatusInternalServerError, err)
			return
		}

		employee := models.Employee{
			ID:       response.Employee.Id,
			Username: response.Employee.Username,
		}

		ctx.Set("Employee", employee)
		ctx.Next()
	}
}

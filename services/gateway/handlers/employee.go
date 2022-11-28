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

type Employee struct {
	employeeClient pb.EmployeeServiceClient
	redisClient    *redis.Client
}

func NewEmployee(
	employeeClient pb.EmployeeServiceClient,
	redisClient *redis.Client,
) Employee {
	return Employee{employeeClient, redisClient}
}

func (e Employee) SignIn(ctx *gin.Context) {
	var req struct {
		Username string
		Password string
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.HttpError(ctx, http.StatusBadRequest, err)
		return
	}

	response, err := e.employeeClient.AuthenticateEmployee(ctx.Request.Context(), &pb.AuthenticateEmployeeRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		helpers.HttpError(ctx, http.StatusUnauthorized, err)
		return
	}

	employee := models.Employee{
		ID:       response.Employee.Id,
		Username: response.Employee.Username,
	}

	token, err := helpers.CreateEmployeeToken(&employee)
	if err != nil {
		helpers.HttpError(ctx, http.StatusInternalServerError, err)
		return
	}

	e.redisClient.Set(ctx.Request.Context(), token, employee.ID, 1*time.Hour)

	helpers.HttpOK(ctx, gin.H{
		"Employee": employee,
		"Token":    token,
	})
}

func (e Employee) SignUp(ctx *gin.Context) {
	var req struct {
		Username string
		Password string
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.HttpError(ctx, http.StatusBadRequest, err)
		return
	}

	response, err := e.employeeClient.CreateEmployee(ctx.Request.Context(), &pb.CreateEmployeeRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		helpers.HttpError(ctx, http.StatusInternalServerError, err)
		return
	}

	employee := models.Employee{
		ID:       response.Employee.Id,
		Username: response.Employee.Username,
	}

	helpers.HttpOK(ctx, gin.H{
		"Employee": employee,
	})
}

func (e Employee) GetAllEmployees(ctx *gin.Context) {
	response, err := e.employeeClient.GetEmployees(ctx, &pb.GetEmployeesRequest{})
	if err != nil {
		helpers.HttpError(ctx, http.StatusInternalServerError, err)
		return
	}

	employees := []models.Employee{}
	for _, protoEmployee := range response.Employees {
		employees = append(employees, models.Employee{
			ID:       protoEmployee.Id,
			Username: protoEmployee.Username,
		})
	}

	helpers.HttpOK(ctx, gin.H{
		"Employees": employees,
	})
}

func (e Employee) GetEmployee(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		helpers.HttpError(ctx, http.StatusBadRequest, err)
		return
	}

	response, err := e.employeeClient.GetEmployee(ctx, &pb.GetEmployeeRequest{
		Id: uint32(id),
	})
	if err != nil {
		helpers.HttpError(ctx, http.StatusBadRequest, err)
	}

	// TODO: Get history booking for employee.

	helpers.HttpOK(ctx, gin.H{
		"Employee": models.Employee{
			ID:       response.Employee.Id,
			Username: response.Employee.Username,
		},
	})
}

func (e Employee) UpdateEmployee(ctx *gin.Context) {
	var req struct {
		Username string `json:"Username,omitempty"`
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

	grpcRequest := &pb.UpdateEmployeeRequest{}
	grpcRequest.Id = uint32(id)
	if req.Username != "" {
		grpcRequest.Username = &req.Username
	}

	response, err := e.employeeClient.UpdateEmployee(ctx, grpcRequest)
	if err != nil {
		helpers.HttpError(ctx, http.StatusInternalServerError, err)
		return
	}

	helpers.HttpOK(ctx, gin.H{
		"Employee": models.Employee{
			ID:       response.Employee.Id,
			Username: response.Employee.Username,
		},
	})
}

func (e Employee) DeleteEmployee(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		helpers.HttpError(ctx, http.StatusBadRequest, err)
		return
	}

	if _, err = e.employeeClient.DeleteEmployee(ctx, &pb.DeleteEmployeeRequest{Id: uint32(id)}); err != nil {
		helpers.HttpError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.AbortWithStatus(http.StatusNoContent)
}

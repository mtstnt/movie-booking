package handlers

import (
	"movie/gateway/pb"

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

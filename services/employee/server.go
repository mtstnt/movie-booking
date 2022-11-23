package main

import (
	"context"
	"movie/employee/pb"

	"golang.org/x/crypto/bcrypt"
)

type Server struct {
	pb.UnimplementedEmployeeServiceServer

	employeeService    Service
	employeeRepository Repository
}

func NewServer(employeeService Service, employeeRepository Repository) Server {
	return Server{
		employeeService:    employeeService,
		employeeRepository: employeeRepository,
	}
}

func (s Server) GetEmployees(ctx context.Context, request *pb.GetEmployeesRequest) (*pb.GetEmployeesResponse, error) {
	employees, err := s.employeeRepository.GetEmployees()
	if err != nil {
		return nil, err
	}

	protoEmployees := make([]*pb.Employee, 0)
	for _, employee := range employees {
		protoEmployees = append(protoEmployees, &pb.Employee{
			Id:       uint32(employee.ID),
			Username: employee.Username,
		})
	}

	return &pb.GetEmployeesResponse{
		Employees: protoEmployees,
	}, nil
}

func (s Server) GetEmployee(ctx context.Context, request *pb.GetEmployeeRequest) (*pb.GetEmployeeResponse, error) {
	id := request.Id

	employee, err := s.employeeRepository.GetEmployee(id)
	if err != nil {
		return nil, err
	}

	return &pb.GetEmployeeResponse{
		Employee: &pb.Employee{
			Id:       uint32(employee.ID),
			Username: employee.Username,
		},
	}, nil
}

func (s Server) CreateEmployee(ctx context.Context, request *pb.CreateEmployeeRequest) (*pb.CreateEmployeeResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	employee := Employee{
		Username: request.Username,
		Password: string(hashedPassword),
	}

	if err := s.employeeRepository.CreateEmployee(&employee); err != nil {
		return nil, err
	}

	return &pb.CreateEmployeeResponse{
		Employee: &pb.Employee{
			Id:       uint32(employee.ID),
			Username: employee.Username,
			Password: string(hashedPassword),
		},
	}, nil
}

func (s Server) UpdateEmployee(ctx context.Context, request *pb.UpdateEmployeeRequest) (*pb.UpdateEmployeeResponse, error) {
	employee, err := s.employeeRepository.GetEmployee(request.Id)
	if err != nil {
		return nil, err
	}

	if request.Username != nil {
		employee.Username = *request.Username
	}

	if err := s.employeeRepository.UpdateEmployee(&employee); err != nil {
		return nil, err
	}

	return &pb.UpdateEmployeeResponse{
		Employee: &pb.Employee{
			Id:       uint32(employee.ID),
			Username: employee.Username,
		},
	}, nil
}

func (s Server) DeleteEmployee(ctx context.Context, request *pb.DeleteEmployeeRequest) (*pb.DeleteEmployeeResponse, error) {
	employee, err := s.employeeRepository.GetEmployee(request.Id)
	if err != nil {
		return nil, err
	}

	if err := s.employeeRepository.DeleteEmployee(&employee); err != nil {
		return nil, err
	}

	return &pb.DeleteEmployeeResponse{}, nil
}

func (s Server) AuthenticateEmployee(ctx context.Context, request *pb.AuthenticateEmployeeRequest) (*pb.AuthenticateEmployeeResponse, error) {
	employee, err := s.employeeService.AuthenticateEmployee(request.Username, request.Password)
	if err != nil {
		return nil, err
	}

	return &pb.AuthenticateEmployeeResponse{
		Employee: &pb.Employee{
			Id:       uint32(employee.ID),
			Username: employee.Username,
		},
	}, nil
}

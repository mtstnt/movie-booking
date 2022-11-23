package main

import (
	"context"
	"movie/user/pb"

	"golang.org/x/crypto/bcrypt"
)

type Server struct {
	pb.UnimplementedUserServiceServer

	userService    Service
	userRepository Repository
}

func NewServer(userService Service, userRepository Repository) Server {
	return Server{
		userService:    userService,
		userRepository: userRepository,
	}
}

func (s Server) GetUsers(ctx context.Context, request *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	users, err := s.userRepository.GetUsers()
	if err != nil {
		return nil, err
	}

	protoUsers := make([]*pb.User, 0)
	for _, user := range users {
		protoUsers = append(protoUsers, &pb.User{
			Id:    uint32(user.ID),
			Email: user.Email,
			Name:  user.Name,
		})
	}

	return &pb.GetUsersResponse{
		Users: protoUsers,
	}, nil
}

func (s Server) GetUser(ctx context.Context, request *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	id := request.Id

	user, err := s.userRepository.GetUser(id)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserResponse{
		User: &pb.User{
			Id:    uint32(user.ID),
			Email: user.Email,
			Name:  user.Name,
		},
	}, nil
}

func (s Server) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := User{
		Email:    request.Email,
		Name:     request.Name,
		Password: string(hashedPassword),
	}

	if err := s.userRepository.CreateUser(&user); err != nil {
		return nil, err
	}

	return &pb.CreateUserResponse{
		User: &pb.User{
			Id:       uint32(user.ID),
			Email:    user.Email,
			Name:     user.Name,
			Password: string(hashedPassword),
		},
	}, nil
}

func (s Server) UpdateUser(ctx context.Context, request *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	user, err := s.userRepository.GetUser(request.Id)
	if err != nil {
		return nil, err
	}

	if request.Email != nil {
		user.Email = *request.Email
	}
	if request.Name != nil {
		user.Name = *request.Name
	}

	if err := s.userRepository.UpdateUser(&user); err != nil {
		return nil, err
	}

	return &pb.UpdateUserResponse{
		User: &pb.User{
			Id:    uint32(user.ID),
			Email: user.Email,
			Name:  user.Name,
		},
	}, nil
}

func (s Server) DeleteUser(ctx context.Context, request *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	user, err := s.userRepository.GetUser(request.Id)
	if err != nil {
		return nil, err
	}

	if err := s.userRepository.DeleteUser(&user); err != nil {
		return nil, err
	}

	return &pb.DeleteUserResponse{}, nil
}

func (s Server) AuthenticateUser(ctx context.Context, request *pb.AuthenticateUserRequest) (*pb.AuthenticateUserResponse, error) {
	user, err := s.userService.AuthenticateEmployee(request.Email, request.Password)
	if err != nil {
		return nil, err
	}

	return &pb.AuthenticateUserResponse{
		User: &pb.User{
			Id:    uint32(user.ID),
			Email: user.Email,
			Name:  user.Name,
		},
	}, nil
}

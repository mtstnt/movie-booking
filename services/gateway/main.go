package main

import (
	"fmt"
	"log"
	"movie/gateway/filters"
	"movie/gateway/handlers"
	"movie/gateway/pb"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcConnMap map[string]*grpc.ClientConn

var (
	ServiceHosts = []string{
		"users",
		"movies",
		"employees",
		"bookings",
		// "notifications", -> Still todo.
	}
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conns, err := dialConns(ServiceHosts, 8000, opts)
	if err != nil {
		return err
	}

	userClient := pb.NewUserServiceClient(conns["users"])
	movieClient := pb.NewMovieServiceClient(conns["movies"])
	employeeClient := pb.NewEmployeeServiceClient(conns["employees"])
	bookingClient := pb.NewBookingServiceClient(conns["bookings"])

	userHandler := handlers.NewUser(userClient, redisClient)
	movieHandler := handlers.NewMovie(employeeClient, movieClient)
	employeeHandler := handlers.NewEmployee(employeeClient, redisClient)
	bookingHandler := handlers.NewBooking(bookingClient)

	userAuthenticatorMiddleware := filters.AuthenticatedAsUser(userClient, redisClient)
	employeeAuthenticatorMiddleware := filters.AuthenticatedAsEmployee(employeeClient)

	router := gin.Default()
	rg := router.Group("/api/v1")

	{
		rg.Group("/users").
			POST("/signin", userHandler.SignIn).
			POST("/signup", userHandler.SignUp)

		rg.Group("/users").
			Use(employeeAuthenticatorMiddleware).
			GET("/", userHandler.GetAllUsers)

		rg.Group("/users").
			Use(userAuthenticatorMiddleware).
			GET("/:id", userHandler.GetUser).
			PUT("/:id", userHandler.UpdateUser).
			DELETE("/:id", userHandler.DeleteUser)
	}

	{
		rg.Group("/movies").
			GET("/", movieHandler.GetMovies).
			GET("/:id", movieHandler.GetMovie)

		rg.Group("/movies").
			Use(filters.AuthenticatedAsEmployee(employeeClient)).
			POST("/", movieHandler.CreateMovie).
			PUT("/:id", movieHandler.UpdateMovie).
			DELETE("/:id", movieHandler.DeleteMovie)
	}

	{
		rg.Group("/employees").
			Use(filters.AuthenticatedAsEmployee(employeeClient)).
			GET("/", employeeHandler.GetAllEmployees).
			GET("/:id", employeeHandler.GetEmployee).
			POST("/", employeeHandler.CreateEmployee).
			PUT("/:id", employeeHandler.UpdateEmployee).
			DELETE("/:id", employeeHandler.DeleteEmployee)
	}

	{
		rg.Group("/bookings").
			Use(filters.AuthenticatedAsUser(userClient, redisClient)).
			GET("/", bookingHandler.GetUserBookings).
			GET("/:id", bookingHandler.GetBooking).
			POST("/", bookingHandler.CreateBooking).
			PUT("/:id", bookingHandler.UpdateBooking).
			DELETE("/:id", bookingHandler.CancelBooking)
	}

	return router.Run("0.0.0.0:8000")
}

func dialConns(services []string, port int, opts []grpc.DialOption) (GrpcConnMap, error) {
	clients := make(GrpcConnMap)
	for _, service := range services {
		c, err := grpc.Dial(fmt.Sprintf("%s:%d", service, port), opts...)
		if err != nil {
			return nil, err
		}
		clients[service] = c
	}
	return clients, nil
}

package main

import (
	"log"
	"movie/movie/pb"
	"net"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	listener, err := net.Listen("tcp", "0.0.0.0:8000")
	if err != nil {
		return err
	}

	db, err := gorm.Open(
		postgres.Open("host=movies_db user=movies password=movies dbname=movies port=5432"),
	)
	if err != nil {
		return err
	}

	// Debug
	if err := resetDB(db); err != nil {
		return err
	}
	if err := db.AutoMigrate(Director{}, Actor{}, Movie{}); err != nil {
		return err
	}
	if err := seed(db); err != nil {
		return err
	}

	movieRepository := NewRepository(db)
	movieServer := NewServer(movieRepository)

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	pb.RegisterMovieServiceServer(grpcServer, movieServer)

	log.Println("Running GRPC Server at port 8000")
	return grpcServer.Serve(listener)
}

func resetDB(db *gorm.DB) error {
	return db.Migrator().DropTable("actors", "directors", "movies", "movie_casts")
}

func seed(db *gorm.DB) error {
	if err := db.Exec(`
		INSERT INTO actors (name) VALUES 
		('Actor1'), ('Actor2'), ('Actor3'), ('Actor4'), ('Actor5'), 
		('Actor6'), ('Actor7'), ('Actor8'), ('Actor9'), ('Actor10')
		ON CONFLICT DO NOTHING
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		INSERT INTO directors (name) VALUES 
		('Director1'), ('Director2'), ('Director3'), ('Director4'), ('Director5')
		ON CONFLICT DO NOTHING
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		INSERT INTO movies (title, synopsis, release_date, director_id) VALUES 
		('Movie1', 'Movie1 synopsis', 1665507600, 1),
		('Movie2', 'Movie2 synopsis', 1664902800, 3),
		('Movie3', 'Movie3 synopsis', 1667062800, 4)
		ON CONFLICT DO NOTHING
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		INSERT INTO movie_casts (movie_id, actor_id) VALUES
		(1, 1), (1, 2), (1, 3), 
		(2, 3), (2, 5), (2, 6),
		(3, 9) 
		ON CONFLICT DO NOTHING
	`).Error; err != nil {
		return err
	}

	return nil
}

package server

import (
	"context"
	"fmt"
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	postgresdriver "github.com/kolya59/easy_normalization/pkg/postgres-driver"
	pb "github.com/kolya59/easy_normalization/proto"
)

type server struct {
	pb.UnimplementedStudentSaverServer
}

func (s *server) SaveStudents(ctx context.Context, in *pb.SaveRequest) (*pb.SaveReply, error) {
	// Convert data
	Students := make([]models.Student, len(in.Students))
	for i, Student := range in.Students {
		Students[i] = *Student
	}

	// Send data in DB
	if err := postgresdriver.SaveStudents(Students); err != nil {
		log.Error().Err(err).Msg("Could not send Students to DB")
		return &pb.SaveReply{}, err
	}

	return &pb.SaveReply{Message: "All is ok"}, nil
}

func StartServer(host, port string) {
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", host, port))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen")
	}
	s := grpc.NewServer()
	pb.RegisterStudentSaverServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("Failed to serve")
	}
}

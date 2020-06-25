package client

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	pb "github.com/kolya59/easy_normalization/proto"
)

func SendStudents(Students []models.Student, host, port string) error {
	// Set up a connection to the server.
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", host, port), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect")
		return fmt.Errorf("failed to connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewStudentSaverClient(conn)

	// Convert Students
	convertedStudents := make([]*models.Student, len(Students))
	for i, Student := range Students {
		convertedStudents[i] = &Student
	}

	// Save Students
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SaveStudents(ctx, &pb.SaveRequest{Students: convertedStudents})
	if err != nil {
		log.Error().Err(err).Msg("Failed to save Students")
		return fmt.Errorf("failed to save Students: %v", err)
	}
	log.Printf("Result: %s", r.GetMessage())
	return nil
}

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	pb "spacegineers_context/protobufs"
	"time"

	"google.golang.org/grpc"
)

const (
	port = 8080
)

type SpaceAndTime struct {
	pb.UnimplementedServicesServer
}

func (s *SpaceAndTime) ContextService(ctx context.Context, req *pb.DamageRequest) (*pb.ResultResponse, error) {
	log.Printf("Received a damage: %v", req.Damage)
	return &pb.ResultResponse{Result: 555}, nil
}

func (s *SpaceAndTime) StreamData(req *pb.NullValue, stream pb.Services_StreamDataServer) error {
	for i := 1; i <= 5; i++ {
		resp := &pb.ResultResponse{Result: int32(i)}

		if err := stream.Send(resp); err != nil {
			return err
		}

		// Simulate some processing time (in a real scenario, perform meaningful work here)
		time.Sleep(time.Second * 2)
	}
	return nil
}

func main() {
	log.Println("Initializing the Context")
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("Failed to listen to the port %d: %v", port, err)
	}
	log.Printf("Listening to the localhost:%d", port)
	grpcServer := grpc.NewServer()
	pb.RegisterServicesServer(grpcServer, &SpaceAndTime{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
	log.Println("Starting the gRPC Server")
}
